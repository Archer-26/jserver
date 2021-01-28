package handler

import (
	"root/pkg/actor"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
)

type IslandHandler struct {
	game interfaces.IGame
}

func InitIslandHandler(g interfaces.IGame) {
	islandHandler := &IslandHandler{
		game: g,
	}
	g.RegistMsg(message.MSG_PLACE_BUILDING_REQ.Int32(), islandHandler.MSG_PLACE_BUILDING_REQ)     // 放置建筑
	g.RegistMsg(message.MSG_CHECKPOINT_FINISH_REQ.Int32(), islandHandler.MSG_ENDPOINT_FINISH_REQ) // 关卡完成
	g.RegistMsg(message.MSG_SAVE_ISLAND_INFO_REQ.Int32(), islandHandler.MSG_SAVE_ISLAND_INFO_REQ) //保存岛屿信息
}

// 放置建筑
func (this *IslandHandler) MSG_PLACE_BUILDING_REQ(actorMsg *actor.ActorMessage) {
	req := actorMsg.Proto().(*message.PlaceBuildingReq)
	role := this.game.PlayerMgr().GetRoleBySession(actorMsg.GateSession())
	role.Island().Place(req.GetLand())

	res := network.NewPbMessage(&message.PlaceBuildingRes{Result: message.MSG_RESULT_SUCCESS}, message.MSG_PLACE_BUILDING_RES.Int32())
	role.SendMsg(res)
}

// 关卡完成
func (this *IslandHandler) MSG_ENDPOINT_FINISH_REQ(actorMsg *actor.ActorMessage) {
	req := actorMsg.Proto().(*message.CheckpointFinishReq)
	role := this.game.PlayerMgr().GetRoleBySession(actorMsg.GateSession())
	role.Island().CheckpointFinish(req)

	res := network.NewPbMessage(&message.CheckpointFinishRes{Result: message.MSG_RESULT_SUCCESS}, message.MSG_CHECKPOINT_FINISH_RES.Int32())
	role.SendMsg(res)
}

// 保存岛屿信息
func (this *IslandHandler) MSG_SAVE_ISLAND_INFO_REQ(actorMsg *actor.ActorMessage) {
	req := actorMsg.Proto().(*message.SaveIslandInfoReq)
	role := this.game.PlayerMgr().GetRoleBySession(actorMsg.GateSession())
	role.Island().SaveJsonInfo(req)

	res := network.NewPbMessage(&message.SaveIslandInfoRes{Result: message.MSG_RESULT_SUCCESS}, message.MSG_SAVE_ISLAND_INFO_RES.Int32())
	role.SendMsg(res)
}
