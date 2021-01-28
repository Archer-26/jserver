package handler_inner

import (
	"root/pkg/actor"
	"root/servers/game/interfaces"
	"root/servers/internal/inner_message/inner"
)

type InnerMsgHandler struct {
	game interfaces.IGame
}

func InitHandler(g interfaces.IGame) {
	innerMsgHandler := &InnerMsgHandler{
		game: g,
	}
	g.RegistMsg(inner.INNER_MSG_GT2G_SESSION_CLOSED.Int32(), innerMsgHandler.INNER_MSG_G2GT_SESSION_CLOSED)
	g.RegistMsg(inner.INNER_MSG_D2G_ROLE_RES.Int32(), innerMsgHandler.INNER_MSG_D2G_ROLE_RES)
}

// gate通知gamesession关闭
func (this *InnerMsgHandler) INNER_MSG_G2GT_SESSION_CLOSED(actorMsg *actor.ActorMessage) {
	recvdata := actorMsg.Proto().(*inner.GT2GSessionClosed)
	this.game.PlayerMgr().SetOfflineRole(recvdata.GateSession)
}

// db返回角色信息
func (this *InnerMsgHandler) INNER_MSG_D2G_ROLE_RES(actorMsg *actor.ActorMessage) {
	recvdata := actorMsg.Proto().(*inner.D2GRoleRes)
	var mod map[string][]byte
	if recvdata.Valid == 0 {
		mod = recvdata.WholeInfo
	} else {
		mod = nil
	}
	// 找不到角色 mod == nil
	this.game.PlayerMgr().CallbackFun(recvdata.CallbackId, recvdata.RID, mod)
}
