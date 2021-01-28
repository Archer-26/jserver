package actor

import (
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/protobuf/proto"
	"root/pkg/abtime"
	"root/pkg/actor/internal/actorpb/protofile"
	"root/pkg/actor/internal/etcd"
	"root/pkg/expect"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/pkg/tools"
	"time"
)

const HEATBEAT_TIME = 20000 //毫秒
type msgHandler func(s *actorSessionHandler, msg proto.Message)
type clientCon struct {
	heatbeatTime int64
	sessionId    uint32
}
type (
	// 管理所有远端的session
	RemoteDispathMgr struct {
		sys      *ActorSystem
		listener network.INetListener
		sfun     chan func() // 所有session相关操作，在RemoteDispathMgr串行执行
		exist    chan bool

		actor2session map[string]uint32               // [actorName]sessionId	 多个远端actor可能映射到同一个sessionId
		addr2session  map[string]*clientCon           // [addr]clientCon 	     主动连接对象，用于管理心跳、防止相同地址的远端重复连接
		sessions      map[uint32]*actorSessionHandler // [sessionID]*INetSession 被动连接对象
		Coder         ICoder

		msgHandler map[int32]msgHandler
		typeMap    *network.MsgTypeMap
		remoteMsg  chan<- IMessage

		etcdEvents <-chan etcd.WatchEvent // 监听远端actor的事件
		newConn    chan network.INetSession
	}

	RemoteActorInfo struct {
		Id   string
		Addr string
	}

	ICoder interface {
		Encode(interface{}) []byte
		Decode([]byte) interface{}
	}
	actorSessionHandler struct {
		*sessionExt
		network.BaseNetHandler
	}
	sessionExt struct {
		Addr string // 如果不为""说明此session是client
		*RemoteDispathMgr
	}
)

func NewRemoteDispathMgr(as *ActorSystem, addr string, etcdEvts <-chan etcd.WatchEvent, in chan<- IMessage, newc chan network.INetSession, coder ICoder) *RemoteDispathMgr {
	if etcdEvts == nil {
		return nil
	}
	mgr := &RemoteDispathMgr{
		sys:           as,
		sessions:      make(map[uint32]*actorSessionHandler),
		actor2session: make(map[string]uint32),
		addr2session:  make(map[string]*clientCon),
		msgHandler:    make(map[int32]msgHandler),
		sfun:          make(chan func(), 5000),
		exist:         make(chan bool),
		etcdEvents:    etcdEvts,
		remoteMsg:     in,
		newConn:       newc,
		Coder:         coder,
	}
	mgr.addr2session[addr] = &clientCon{} // 因为actor会在最后初始化，所以etcd会推送自己本端的actor过来，这里防止自己连自己
	mgr.typeMap = network.NewMsgTypeMap()
	mgr.typeMap.InitMsgParser("protofile", "MSG_TYPE")
	mgr.listener = network.StartTcpListen(addr, &network.StreamCodec{}, &actorSessionHandler{}, network.SetSSExtra(&sessionExt{Addr: "", RemoteDispathMgr: mgr}))

	mgr.registMsgHandler(int32(protofile.MSG_TYPE_ACTOR_HEATBEAT_PING), heatbeatpingHander)
	mgr.registMsgHandler(int32(protofile.MSG_TYPE_ACTOR_HEATBEAT_PONG), heatbeatpongHander)
	mgr.registMsgHandler(int32(protofile.MSG_TYPE_ACTOR_REGIST), registHander)
	mgr.registMsgHandler(int32(protofile.MSG_TYPE_ACTOR_MESSAGE), recvHander)
	return mgr
}

func (this *RemoteDispathMgr) Startup() *RemoteDispathMgr {
	go func() {
		heatbeatTimer := time.NewTicker(HEATBEAT_TIME * time.Millisecond)
		leastBeatTime := int64(0)
		defer func() {
			heatbeatTimer.Stop()
			this.listener.Stop()
		}()

		for {
			select {
			case f := <-this.sfun:
				tools.Try(func() {
					f()
				}, nil)
			case ev := <-this.etcdEvents:
				tools.Try(func() {
					for actorId, val := range ev {
						if val.T == mvccpb.PUT {
							actorInfo := &RemoteActorInfo{}
							err := json.Unmarshal([]byte(val.Evt), actorInfo)
							expect.Nil(err,log.Fields{"error": err, "Evt": val.Evt})

							sessionId := uint32(0)
							if sid, ok := this.actor2session[actorId]; ok {
								// 有一个actor重复了，说明对应的远端出问题，把此actor的session关联的所有actor清除,等待对方重新put
								log.KVs(log.Fields{"actorId": actorId, "sid": sid, "val": val}).Warn(colorized.Red("remote actor already exists"))
								this.clearSession(sid)
								sessionId = this.newConnect(actorInfo.Addr, actorInfo.Id)
							} else if cli, ok := this.addr2session[actorInfo.Addr]; ok {
								if cli.sessionId == 0 { // 是自己，不需要连接
									//log.WarnOld("%v %v is oneself!", actorId, actorInfo.Addr)
									break
								}
								// 同一个远端的不同actor
								sessionId = cli.sessionId
							} else {
								// 正常情况新连接都直接走这里
								sessionId = this.newConnect(actorInfo.Addr, actorInfo.Id)
								log.KVs(log.Fields{"actor": actorInfo.Id, "addr": actorInfo.Addr}).Info(colorized.Cyan("etcd a new client session succes!"))
							}
							this.setActor2Session(actorInfo.Id, sessionId)
						} else if val.T == mvccpb.DELETE {
							sessionId, e := this.actor2session[actorId]
							if !e {
								break
							}
							if sessionId != 0 && this.sessions[sessionId] != nil { // etcd DELETE
								this.sessions[sessionId].INetSession.Stop()
								log.KV("actorId", actorId).KV("sessionId", sessionId).Info("etcd DELETE!")
							}
						} else {
							log.KV("val.T", val.T).Error("etcdEvents error")
						}
					}
				}, nil)
			case <-heatbeatTimer.C:
				tools.Try(func() {
					for addr, cli := range this.addr2session {
						session := this.sessions[cli.sessionId] // 心跳
						if session == nil {
							if cli.sessionId != 0 { // 0 为本地特殊标记
								log.KV("addr", addr).Info(" session alreay closed heatbeat delete")
								delete(this.addr2session, addr)
							}
							continue
						}
						// 心跳超时处理
						if cli.heatbeatTime < leastBeatTime {
							log.KVs(log.Fields{
								"heatbeatTime":  cli.heatbeatTime,
								"leastBeatTime": leastBeatTime,
								"sessionId":     cli.sessionId,
							}).Warn("heat timeout")
							session.INetSession.Stop()
							// todo...依赖etcd 不需要重连机制？
						} else {
							// 发送心跳
							netMsg := network.NewPbMessage(&protofile.ActorHeatbeatPing{Addr: addr}, int32(protofile.MSG_TYPE_ACTOR_HEATBEAT_PING))
							session.SendMsg(netMsg.Buffer())
						}
					}
					leastBeatTime = abtime.Milliseconds()
				}, nil)
			case <-this.exist:
				return
			}
		}
	}()
	return this
}

func (this *RemoteDispathMgr) Stop() {
	close(this.exist)
}

// 远端actor发送消息
func (this *RemoteDispathMgr) Send(msg *ActorMessage) {
	this.sfun <- func() { // Actor发消息给远端
		fields := log.Fields{"source": msg.SourceId(), "target": msg.TargetId()}
		sessionId, ok := this.actor2session[msg.TargetId()]
		if !ok {
			log.KVs(fields).Red().Error("not found remote actor")
			this.sys.event._dispatchEvent(Ev_remoteInexistent{msg.TargetId()})
			return
		}
		s := this.sessions[sessionId]
		if s == nil {
			log.KVs(fields).KV("sessionId", sessionId).Red().Error("not found session")
			this.sys.event._dispatchEvent(Ev_remoteInexistent{msg.TargetId()})
			return
		}
		pbMessage := &protofile.ActorMessage{
			GateSession: msg.gateSession,
			SourceId:    msg.sourceId,
			TargetId:    msg.targetId,
			MsgId:       msg.msgId,
			Data:        s.Coder.Encode(msg.data),
		}
		netMsg := network.NewPbMessage(pbMessage, int32(protofile.MSG_TYPE_ACTOR_MESSAGE))
		buffer := netMsg.Buffer()
		s.SendMsg(buffer)
		this.sys.msgPool.Put(msg)
	}
}

func (this *RemoteDispathMgr) newConnect(addr, actorId string) uint32 {
	fields := log.Fields{"actor": actorId, "addr": addr}
	// 先建立连接，再添加actor
	cli, ok := network.StartTcpClient(addr, &network.StreamCodec{}, &actorSessionHandler{}, network.SetCSExtra(&sessionExt{Addr: addr, RemoteDispathMgr: this}))
	if !ok {
		log.KVs(fields).Error("network.StartTcpClient faied")
		return 0
	}
	this.newConn <- cli.NetSession() // 连接成功，本地所有actor发送给cli
	this.addr2session[addr] = &clientCon{
		sessionId:    cli.NetSession().Id(),
		heatbeatTime: abtime.Milliseconds(),
	}
	log.KVs(fields).Cyan().Info("RemoteDispathMgr newConnect!")
	return cli.NetSession().Id()

}

func (this *RemoteDispathMgr) clearSession(sessionId uint32) {
	for id, sesid := range this.actor2session {
		if sessionId == sesid {
			delete(this.actor2session, id)
		}
	}
	for key, cli := range this.addr2session {
		if sessionId == cli.sessionId {
			delete(this.addr2session, key)
		}
	}
	if s := this.sessions[sessionId]; s != nil {
		s.INetSession.Stop()
		delete(this.sessions, sessionId)
	}
}
func (this *RemoteDispathMgr) registMsgHandler(msgId int32, fn msgHandler) {
	if _, ok := this.msgHandler[msgId]; ok {
		log.KV("msgID", msgId).ErrorStack(2, "registMsgHandler existence")
		return
	}
	this.msgHandler[msgId] = fn
}

func (this *RemoteDispathMgr) setActor2Session(actorId string, sessionId uint32) {
	//if _,existence := this.actor2session[actorId];!existence{
	//	 todo...通知dispatcher激活actor
	//}
	this.actor2session[actorId] = sessionId
	this.sys.event._dispatchEvent(Ev_newActor{ActorId: actorId, Remote: REMOTE})
}

///////////////////////////////////////// actorSessionHandler /////////////////////////////////////////////
func (this *actorSessionHandler) OnSessionCreated() {
	ext, ok := this.Load("ext")
	if !ok {
		log.KV("sessionId", this.Id()).Error("OnSessionCreated, can't find ext")
		return
	}
	this.sessionExt = ext.(*sessionExt)
	// 无论是主动连接还是监听新连接，都会走这里
	this.sfun <- func() { // 新连接
		this.sessions[this.Id()] = this
	}
}

func (this *actorSessionHandler) OnSessionClosed() {
	this.sfun <- func() { // 连接关闭
		delete(this.sessions, this.Id())
		if this.Addr != "" {
			delete(this.addr2session, this.Addr)
		}
		log.KV("sessionId", this.Id()).KV("Addr", this.Addr).Info(colorized.Gray("OnSessionClosed remove Closed"))
	}
}

func (this *actorSessionHandler) OnRecv(bytes []byte) {
	fields := log.Fields{"sessionId": this.Id(), "addr": this.RemoteAddr(), "len(data)": len(bytes)}
	netMsg := network.NewBytesMessageParse(bytes, this.typeMap)
	if netMsg == nil {
		log.KVs(fields).Error("bad client netMsg")
		this.INetSession.Stop()
		return
	}
	this.sfun <- func() { // session收到数据
		if onHander, ok := this.msgHandler[netMsg.MsgId()]; ok {
			actorMsg := netMsg.Proto()
			tools.Try(func() { onHander(this, actorMsg) }, nil)
		} else {
			log.KVs(fields).KV("msgId", netMsg.MsgId()).Error("HandleMsg msgId not found")
		}
	}
}
