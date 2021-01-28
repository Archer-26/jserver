package handler

import (
	"root/pkg/actor"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
)

type DiceHandler struct {
	game interfaces.IGame
}

func InitDiceHandler(g interfaces.IGame) {
	diceHandler := &DiceHandler{
		game: g,
	}
	g.RegistMsg(message.MSG_DICE_SAVE_INFO_REQ.Int32(), diceHandler.MSG_DICE_SAVE_INFO_REQ) // 保存游戏信息
}

// 保存游戏信息
func (this *DiceHandler) MSG_DICE_SAVE_INFO_REQ(actorMsg *actor.ActorMessage) {
	req := actorMsg.Proto().(*message.DiceSaveInfoReq)
	role := this.game.PlayerMgr().GetRoleBySession(actorMsg.GateSession())
	pb := role.Dice().SaveGameInfo(req)
	res := network.NewPbMessage(pb, message.MSG_DICE_SAVE_INFO_RES.Int32())
	role.SendMsg(res)
}
