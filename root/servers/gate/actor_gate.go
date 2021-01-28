package gate

import (
	"root/internal/common"
	"root/pkg/actor"
	"root/pkg/expect"
	"root/pkg/iniconfig"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/message"
	"time"
)

/*
 * 针对客户端连接的网关, 负责逻辑actor和client之间的消息中转，不处理任何消息
 */
type (
	Gate struct {
		Addr  string
		Parse bool

		actor.IActor
		sessions map[uint32]*UserSessionHandler

		listener network.INetListener // 用户连接端口
		typemap  *network.MsgTypeMap
	}
)

func (this *Gate) Init(a actor.IActor) {
	this.IActor = a
	this.sessions = make(map[uint32]*UserSessionHandler)
	this.Addr = iniconfig.String("gate_addr")
	expect.True(this.Addr != "")

	this.listener = network.StartTcpListen(this.Addr, &network.StreamCodec{}, &UserSessionHandler{}, network.SetSSMaxRead(1024), network.SetSSExtra(this))
	this.typemap = network.NewMsgTypeMap().InitMsgParser("message", "MSG")
	this.RegistEvent(actor.AEV_REMOTE_INEXISTENT, this.OnRemoteInexistent)

	this.AddTimer(time.Hour, -1, this.checkDeadSession)
}

func (this *Gate) Stop() bool {
	this.CancelEvent(actor.AEV_REMOTE_INEXISTENT)
	this.listener.Stop()
	return true
}

// 定期检查并清理死链接
func (this *Gate) checkDeadSession(dt int64) {
	for id, session := range this.sessions {
		if session.Alive {
			session.Alive = false
		} else {
			session.Stop()
			delete(this.sessions, id)
		}
	}
}

// 处理actorcore抛来的事件
func (this *Gate) OnRemoteInexistent(event actor.ActorEvent) {
	switch event.EType() {
	case actor.AEV_REMOTE_INEXISTENT:
		evData, ok := event.(actor.Ev_remoteInexistent)
		expect.True(ok,log.Fields{"etype": event.EType()})

		log.KV("remote actor", evData.ActorId).Warn("remote actor is inexistent")
	}
}

// 所有消息，直接转发给用户
func (this *Gate) HandleMessage(actorMsg *actor.ActorMessage) {
	msgId := actorMsg.MsgId()
	if inner.INNER_MSG_INNER_SEGMENT_BEGIN.Int32() < msgId && msgId < inner.INNER_MSG_INNER_SEGMENT_END.Int32() {
		this.InnerHandler(actorMsg) // 内部消息，单独处理
		return
	}

	// 用户消息，直接转发给用户
	actorId, sessionId := common.SplitGateSession(actorMsg.GateSession())
	logInfo := log.Fields{
		"own":         this.GetID(),
		"gateSession": actorMsg.GateSession(),
		"sourceId":    actorMsg.SourceId(),
		"msgId":       actorMsg.MsgId(),
		"msgName":     message.MSG(actorMsg.MsgId()),
	}
	expect.True(this.GetID() == actorId,logInfo)

	userSessionHandler := this.sessions[sessionId]
	if userSessionHandler == nil {
		log.KVs(logInfo).Warn("cannot find sessionId")
		return
	}
	log.KVs(logInfo).Info("server message to user")
	userSessionHandler.SendMsg(actorMsg.NetMessage().Buffer())
	return
}
