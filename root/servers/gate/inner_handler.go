package gate

import (
	"root/internal/common"
	"root/pkg/actor"
	"root/pkg/expect"
	"root/pkg/log"
	"root/servers/internal/inner_message/inner"
)

func (this *Gate) InnerHandler(actorMsg *actor.ActorMessage) bool {
	switch actorMsg.MsgId() {
	case inner.INNER_MSG_L2GT_SESSION_ASSIGN_GAME.Int32(): // login分配游戏服，通知gate绑定用户gameActor
		this.INNER_MSG_L2G_SESSION_ASSIGN_GAME(actorMsg)
	case inner.INNER_MSG_L2GT_USER_SESSION_DISABLED.Int32(): // login通知gate 用户旧session失效
		this.INNER_MSG_L2GT_USER_SESSION_DISABLED(actorMsg)
	default:
		return false
	}
	return true
}

func (this *Gate) INNER_MSG_L2G_SESSION_ASSIGN_GAME(actorMsg *actor.ActorMessage) {
	recvData := actorMsg.Proto().(*inner.L2GTSessionAssignGame)
	gate, sessionId := common.SplitGateSession(actorMsg.GateSession())
	fields := log.Fields{"actor": this.GetID(), "gateSession": actorMsg.GateSession(), "source": actorMsg.SourceId()}
	expect.True(this.GetID() == gate,fields)

	session := this.sessions[sessionId]
	if session == nil {
		log.KVs(fields).Warn("session was closed")
		return
	}
	session.GameActor = recvData.GameActorId
	log.KVs(log.Fields{"session": recvData.GateSession, "game actorId": recvData.GameActorId}).
		Info("INNER_MSG_L2G_SESSION_ASSIGN_GAME")
}

func (this *Gate) INNER_MSG_L2GT_USER_SESSION_DISABLED(actorMsg *actor.ActorMessage) {
	recvData := actorMsg.Proto().(*inner.L2GTUserSessionDisabled)
	gate, sessionId := common.SplitGateSession(recvData.GetGateSession())
	fields := log.Fields{"actor": this.GetID(), "gateSession": actorMsg.GateSession(), "source": actorMsg.SourceId()}
	expect.True(this.GetID() == gate,fields)

	if this.GetID() != gate {
		log.KVs(fields).Error("exception this.GetID() != gate local")
		return
	}
	session := this.sessions[sessionId]
	if session == nil {
		log.KVs(fields).Warn("recvice disabled session msg,but session is nil")
		return
	}
	session.INetSession.Stop()
	log.KVs(fields).Info("INNER_MSG_L2GT_USER_SESSION_DISABLED ")
}
