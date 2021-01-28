package actor

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"root/pkg/abtime"
	"root/pkg/log"
	"root/pkg/network"
	"sync/atomic"
	"time"
)

type (
	IActor interface {
		GetID() string
		MailCount() int
		AddTimer(interval time.Duration, trigger_times int32, callback abtime.FuncCallback) int64
		CancelTimer(timerId int64)
		ResetTimerMgr()
		RegistEvent(eventType ActorEventType, evf Callback)
		CancelEvent(eventType ActorEventType)
		SuspendStop()
	}

	// acotr 处理接口
	IFHandler interface {
		Init(actor IActor)
		Stop() bool // true 立刻停止，false 延迟停止
		HandleMessage(*ActorMessage)
	}

	// actor邮箱消息基类
	IMessage interface {
		SourceId() string
		TargetId() string
	}
)

type (
	ActorOption func(*actor)
	// 运行单元
	actor struct {
		id       string
		handler  IFHandler
		mailBox  chan IMessage // todo 考虑用无锁队列优化
		sys      *ActorSystem
		timerMgr *abtime.TimerMgr
		remote   bool // 是否能被远端发现 默认为true, 如果是本地actor,手动设Localized()后,再注册

		exist chan struct{} // 用于外部通知结束
		stop  int32     // 用于内部自己结束	(一些actor的关闭是由其他actor调用,防止core退出时重复close(exist)的情况,故特意添加此变量,用于主动调用SuspendStop())
	}

	ActorMessage struct {
		sourceId    string
		targetId    string
		gateSession string
		msgId       int32
		data        interface{}
	}

	// 只能用于本地Actor，方便处理闭包
	localMessage struct {
		sourceId   string
		targetId   string
		fname      string
		funHandler func()
	}
)

// 创建actor
// id 		actorId外部定义
// slothandler  消息处理模块
// op  修改默认属性
func NewActor(id string, handler IFHandler, op ...ActorOption) IActor {
	actor := &actor{
		id:       id,
		handler:  handler,
		mailBox:  make(chan IMessage, 1000),
		exist:    make(chan struct{}),
		remote:   true, // 默认都能被远端发现
		timerMgr: abtime.NewTimerMgr(),
	}
	for _, f := range op {
		f(actor)
	}
	return actor
}
func SetMailBoxSize(boxSize int) ActorOption {
	return func(a *actor) {
		a.mailBox = make(chan IMessage, boxSize)
	}
}
func AvailWheel() ActorOption {
	return func(a *actor) {
		a.timerMgr.AvailWheel()
	}
}
func SetLocalized() ActorOption {
	return func(a *actor) {
		a.remote = false
	}
}

// 获取actorID
func (this *actor) GetID() string {
	return this.id
}

func (this *actor) MailCount() int {
	return len(this.mailBox)
}

// 逻辑层主动关闭actor调用此函数
func (this *actor) SuspendStop() {
	atomic.StoreInt32(&this.stop, 1)
}

// Push一个消息
func (this *actor) Push(msg IMessage) {
	if msg == nil {
		return
	}
	l := len(this.mailBox)
	c := cap(this.mailBox)
	if l > c*2/3 {
		log.KVs(log.Fields{
			"actor": this.id,
			"len":   l,
			"cap":   c,
		}).Warn("警告！队列消息超过三分之二了")
	}
	this.mailBox <- msg
}

// 注册事件回调
func (this *actor) RegistEvent(eventType ActorEventType, evf Callback) {
	this.sys.event.registEvent(this.id, eventType, evf)
}

// 注销事件回调
func (this *actor) CancelEvent(eventType ActorEventType) {
	this.sys.event.cancelEvent(this.id, eventType)
}

// 添加计时器,每个actor独立一个计时器
// interval 	 单位nanoseconds
// trigger_times 执行次数 -1 无限次
// callback 	 只能是主线程回调
func (this *actor) AddTimer(interval time.Duration, trigger_times int32, callback abtime.FuncCallback) int64 {
	if this.timerMgr == nil {
		return -1
	}
	now := abtime.Now().UnixNano()
	newTimer, e := abtime.NewTimer(now, now+interval.Nanoseconds(), trigger_times, callback)
	if e != nil {
		log.KV("error", e).ErrorStack(3, "AddTimer failed")
		return -1
	}
	return this.timerMgr.AddTimer(newTimer, false)
}

// 删除一个定时器
func (this *actor) CancelTimer(timerId int64) {
	this.timerMgr.CancelTimer(timerId)
}

// 重置清空
func (this *actor) ResetTimerMgr() {
	this.timerMgr.Reset()
}

/////////////////////////////////////////////// IMessage /////////////////////////////////////////////////////
// ActorMessage
func (this *ActorMessage) MsgId() int32 {
	return this.msgId
}
func (this *ActorMessage) SourceId() string {
	return this.sourceId
}
func (this *ActorMessage) TargetId() string {
	return this.targetId
}
func (this *ActorMessage) GateSession() string {
	return this.gateSession
}
func (this *ActorMessage) Proto() proto.Message {
	return this.data.(*network.Message).Proto()
}
func (this *ActorMessage) NetMessage() *network.Message {
	return this.data.(*network.Message)
}
func (this *ActorMessage) String() string {
	return fmt.Sprintf("msgId:%v source:%v target:%v gateSession:[%v]",
		this.msgId, this.sourceId, this.targetId, this.gateSession)
}

// localMessage
func (this *localMessage) SourceId() string {
	return this.sourceId
}
func (this *localMessage) TargetId() string {
	return this.targetId
}
func (this *localMessage) String() string {
	return fmt.Sprintf("source:%v target:%v funname:%v", this.sourceId, this.targetId, this.fname)
}
