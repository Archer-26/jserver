package actor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"root/pkg/abtime"
	"root/pkg/actor/internal/actorpb/protofile"
	"root/pkg/actor/internal/etcd"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/pkg/tools"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

/* 所有actor的驱动器和调度器*/
const LOCAL = 1  //本地actor
const REMOTE = 2 //远端actor

type SystemOption func(*ActorSystem)
type ActorSystem struct {
	wg        *sync.WaitGroup
	actorAddr string // 远程actor连接端口

	// 本地actor 管理相关
	newList      chan IActor
	closeList    chan string
	actorCache   sync.Map      // 所有本地actor
	msgQueue     chan IMessage // actor的消息分发队列
	msgPool      sync.Pool
	localmsgPool sync.Pool
	exit         chan struct{} // 进程退出
	bexist       bool

	// 远端actor相关
	remoteCoder ICoder
	remoteMgr   *RemoteDispathMgr
	etcdService *etcd.Etcd
	etcdAddr    string // 需要连接的etcd端点
	etcdPrefix  string

	event *EvDispatcher
}

func NewActorSystem(op ...SystemOption) *ActorSystem {
	sys := &ActorSystem{
		wg:         &sync.WaitGroup{},
		newList:    make(chan IActor, 100),
		closeList:  make(chan string, 100),
		msgQueue:   make(chan IMessage, 1000), // 不够再加
		exit:       make(chan struct{}),
		actorAddr:  "",
		etcdAddr:   "",
		etcdPrefix: "",
	}
	sys.event = NewActorEvents(sys)
	sys.msgPool.New = func() interface{} { return &ActorMessage{} }
	sys.localmsgPool.New = func() interface{} { return &localMessage{} }

	for _, f := range op {
		f(sys)
	}
	// 命令行
	sys.RegistCmd("", "actorinfo", sys.actorInfo)
	sys.RegistCmd("", "detect", sys.detect)
	return sys
}

// 阻塞处理所有actor,主线程最后调用
func (s *ActorSystem) Startup() {
	log.Info(colorized.Blue("Startup...."))
	s.wg.Add(1)

	// 启动流程 远端管理器->etcd服务发现->本地actor调度器
	newConn := make(chan network.INetSession)
	event_chan := make(chan etcd.WatchEvent, 1000) // 此chan生产者是etcd 消费者是RemoteDispatherMgr
	s.remoteMgr = NewRemoteDispathMgr(s, s.actorAddr, event_chan, s.msgQueue, newConn, s.remoteCoder).Startup()
	s.etcdService = etcd.NewEtcd(s.etcdAddr, s.etcdPrefix, event_chan).Startup()

	catchs := make(chan os.Signal)
	signal.Notify(catchs, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		select {
		case sig := <-catchs:
			log.KV("signal", sig.String()).Red().Warn("sys signal catched")
			//shutdown() 通知所有actor执行关闭
			s.actorCache.Range(func(key, value interface{}) bool {
				a := value.(*actor)
				a.exist <- struct{}{}
				return true
			})
			close(s.newList)
			s.bexist = true
			s.wg.Done()
			go func() { // 关闭后，core继续运行，等待所有acotr都完成关闭再退出进程
				s.wg.Wait()
				s.remoteMgr.Stop()
				close(s.exit)
			}()
		case v := <-s.newList:
			if !s.bexist { // 进程退出状态，不再处理新的actor
				s._add(v.(*actor))
			}
		case id := <-s.closeList:
			s.actorCache.Delete(id)
		case msg := <-s.msgQueue:
			r := false	// 消息能否发送给远端
			if _,ok := msg.(*ActorMessage);ok{
				r = true
			}
			val, ok := s.actorCache.Load(msg.TargetId())
			if ok {
				target := val.(*actor)
				target.Push(msg)
			} else if r {
				s.remoteMgr.Send(msg.(*ActorMessage)) // 如果本地找不到，发送给网络
			} else {
				log.KV("target", msg.TargetId()).Red().Error(" not found local")
			}
		case session := <-newConn:
			actors := []string{}
			s.actorCache.Range(func(key, value interface{}) bool {
				actorObj := value.(*actor)
				if actorObj.remote {
					actors = append(actors, (key).(string))
				}
				return true
			})
			netMsg := network.NewPbMessage(&protofile.ActorRegist{ActorId: actors}, int32(protofile.MSG_TYPE_ACTOR_REGIST))
			session.SendMsg(netMsg.Buffer())
		case <-s.exit:
			log.Info("shutdown")
			return
		}
	}
}

// 注册actor，外部创建对象，保证ActorId唯一性
func (s *ActorSystem) Regist(actor IActor) {
	s.newList <- actor
}

//
func (s *ActorSystem) LocalActorFuzzy(f func(reg string) bool) []IActor {
	ret := []IActor{}
	s.actorCache.Range(func(key, value interface{}) bool {
		if f(key.(string)) {
			ret = append(ret, value.(IActor))
		}
		return true
	})
	return ret
}

// 启动前设置
func (s *ActorSystem) RemoteCoder(coder ICoder) {
	s.remoteCoder = coder
}

// actor之间发送消息
// sourceid 发送源actor
// targetid 目标actor
// pbMsg 消息内容
func (s *ActorSystem) Send(gateSession, sourceId, targetId string, pbMsg *network.Message) {
	logSend := fmt.Sprintf("%v Send ActorMessage", sourceId)
	log.KVs(log.Fields{"targetId": targetId, "MsgId": pbMsg.MsgId()}).Bule().Info(logSend)
	msg := s.msgPool.Get().(*ActorMessage)
	msg.gateSession = gateSession
	msg.sourceId = sourceId
	msg.targetId = targetId
	msg.data = pbMsg
	msg.msgId = pbMsg.MsgId()
	s.msgQueue <- msg
}

// actor之间发送消息(targetid 必须保证是本地actor)
// sourceid 发送源actor
// targetid 目标actor
// funName  没有业务逻辑，日志记录用,防止闭包导致死循环等问题无法追溯
// f        闭包函数
func (s *ActorSystem) LocalSend(sourceId, targetId, funName string, f func()) {
	logSend := fmt.Sprintf("%v Send localMessage", sourceId)
	log.KVs(log.Fields{"targetId": targetId, "funName": funName}).Bule().Info(logSend)
	msg := s.localmsgPool.Get().(*localMessage)
	msg.sourceId = sourceId
	msg.targetId = targetId
	msg.fname = funName
	msg.funHandler = f
	s.msgQueue <- msg
}

func (s *ActorSystem) _add(actor *actor) {
	if actor == nil {
		log.Error("actor == nil")
		return
	}
	if _, ok := s.actorCache.Load(actor.id); ok {
		log.KV("actor", actor.id).ErrorStack(4, "REPEATED ACTOR !!!!")
		return
	} else {
		s.actorCache.Store(actor.id, actor)
		actor.sys = s
		go _runActor(s, actor)
	}
}

func _runActor(sys *ActorSystem, actor *actor) {
	sys.wg.Add(1)
	// etcd注册actor信息
	if actor.remote {
		info := RemoteActorInfo{Id: actor.id, Addr: sys.actorAddr}
		jsonVal, err := json.Marshal(&info)
		if err != nil {
			log.KVs(log.Fields{"err": err, "info": info}).Error("json marshal error")
			return
		}
		sys.etcdService.RegistEtcd(actor.id, string(jsonVal))
	}

	tools.Try(func() {
		log.KV("actor", actor.GetID()).Yellow().Info("ACTOR startup!")
		actor.handler.Init(actor)
	}, nil)

	defer func() {
		log.KV("actor", actor.GetID()).Red().Info("ACTOR Done!")
		sys.closeList <- actor.GetID()
		sys.wg.Done()
	}()

	sys.event._dispatchEvent(Ev_newActor{ActorId: actor.GetID(), Remote: LOCAL})

	up_timer := time.NewTicker(time.Millisecond * 100)
	// actor逻辑循环
	for {
		select {
		case <-actor.exist: // main线程想退出，逻辑层被动关闭
			actor.exist = nil // set nil 让select不在监听此变量
			if actor.handler.Stop() {return}
		default:
			select {
			case <-actor.exist:
				actor.exist = nil
				if actor.handler.Stop() {return}
			case <-up_timer.C:
				tools.Try(func() {
					actor.timerMgr.Update(abtime.Now().UnixNano())
				}, nil)
			case msg := <-actor.mailBox:
				l,c := len(actor.mailBox),cap(actor.mailBox)
				if l > c*2/3 {
					log.KVs(log.Fields{
						"actor": actor.id,
						"l":     l,
						"c":     c,
					}).Warn(" too much msg watting process in mailBox of actor !!!!!!")
				}

				err := watchWarn(msg, func() {
					// 一次只处理一条消息
					switch message := msg.(type) {
					case *ActorMessage:
						logSend := fmt.Sprintf("%v Recv ActorMessage", actor.GetID())
						log.KVs(log.Fields{"message": message.String()}).Bule().Info(logSend)
						if sys.detectMsg(message.MsgId(), message.TargetId(), message.SourceId()) {
							break
						}
						actor.handler.HandleMessage(message)
						sys.msgPool.Put(message)
					case *localMessage:
						logSend := fmt.Sprintf("%v Recv localMessage", actor.GetID())
						log.KVs(log.Fields{"msg": message.String()}).Bule().Info(logSend)
						message.funHandler()
						sys.localmsgPool.Put(message)
					default:
						log.KVs(log.Fields{"source": message.SourceId(), "target": message.TargetId()}).Red().Error("error actor message type,neither ActorMessage nor localMessage")
					}
				})
				if err != nil {
					log.KV("error", err).Red().Warn("watchWarn Handle CAUTION")
				}
			}
		}

		// 逻辑层调用SuspendStop()
		if atomic.LoadInt32(&actor.stop) == 1 {
			return
		}
	}
}

func watchWarn(msg IMessage, handle func()) error {
	beginTime := abtime.Milliseconds()
	tools.Try(func() { handle() }, nil)
	endTime := abtime.Milliseconds()

	if endTime-beginTime > 2000 {
		var (
			gsess, souid, tarid, fname string
			msgId                      int32
		)
		switch msgData := msg.(type) {
		case *ActorMessage:
			gsess, souid, tarid, msgId = msgData.gateSession, msgData.sourceId, msgData.targetId, msgData.msgId
		case *localMessage:
			souid, tarid, fname = msgData.sourceId, msgData.targetId, msgData.fname
		}

		return errors.New(fmt.Sprintf("handle processed too long process-abtime:%v gateSession:%v sourceId:%v targetId:%v MsgId:%v fname:%v",
			endTime-beginTime, gsess, souid, tarid, msgId, fname))
	}

	return nil
}

func (s *ActorSystem) detectMsg(msgId int32, target, source string) bool {
	if msgId == int32(protofile.MSG_TYPE_ACTOR_DETECT) { // 收到远端的探测信号
		detect := network.NewPbMessage(nil, int32(protofile.MSG_TYPE_ACTOR_DETECT_RESP))
		s.Send("", target, source, detect)
		log.KVs(log.Fields{"source": source, "target": target}).Info("recv detect ok!")
		return true
	} else if msgId == int32(protofile.MSG_TYPE_ACTOR_DETECT_RESP) { // 收到远端回复的探测信号
		log.KVs(log.Fields{"source": source, "target": target}).Info("recv detect resp ok!")
		return true
	}
	return false
}
