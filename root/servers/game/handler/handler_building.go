package handler

import (
	"root/internal/common"
	"root/internal/config/config_global"
	"root/internal/config/config_go"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
)

type BuildingHandler struct {
	game interfaces.IGame
}

func InitBuildingHandler(g interfaces.IGame) {
	handler := &BuildingHandler{
		game: g,
	}
	g.RegistMsg(message.MSG_BUILDING_LEVEL_UP_REQ.Int32(), handler.MSG_BUILDING_LEVEL_UP_REQ)
	g.RegistMsg(message.MSG_BUILDING_STAR_UP_REQ.Int32(), handler.MSG_BUILDING_STAR_UP_REQ)

	g.RegistMsg(message.MSG_BUILDING_PROP_UP_REQ.Int32(), handler.MSG_BUILDING_PROP_UP_REQ)

}

func (this *BuildingHandler) MSG_BUILDING_LEVEL_UP_REQ(actorMsg *actor.ActorMessage) {
	//req := actorMsg.Proto().(*message.BuildingLevelUpReq)
	//role := this.game.PlayerMgr().GetRoleBySession(actorMsg.GateSession())
	//building := role.Building().GetBuilding(req.Id)
	//if building == nil {
	//	this.BuildingLevelUpRes(role.RID(), actorMsg, message.MSG_RESULT_NO_THIS_BUILDING)
	//	return
	//}
	//cfg := config_go.GetBaseBuildingLevel(config_global.GetBuildingLevelKey(req.Id, req.UptoLevel))
	//if cfg == nil { //检测想升到的等级是否存在
	//	this.BuildingLevelUpRes(role.RID(), actorMsg, message.MSG_RESULT_CFG_NO_THIS_PARAM)
	//	return
	//}
	//if req.UptoLevel <= building.Lv {
	//	this.BuildingLevelUpRes(role.RID(), actorMsg, message.MSG_RESULT_CLIENT_WRONG_PARAM)
	//	return
	//}
	//var costItem = make(map[int32]int64)
	//for level := building.Lv; level < req.UptoLevel; level++ {
	//	cfg := config_go.GetBaseBuildingLevel(config_global.GetBuildingLevelKey(req.Id, level))
	//	if cfg == nil {
	//		this.BuildingLevelUpRes(role.RID(), actorMsg, message.MSG_RESULT_CFG_NO_THIS_PARAM)
	//		return
	//	}
	//	for i := 0; i < cfg.Len_LEVEL_UP_ITEM(); i++ {
	//		costItem[cfg.Get_LEVEL_UP_ITEM(i)] += int64(cfg.Get_LEVEL_UP_NUM(i))
	//	}
	//}
	//
	//if !role.ItemsEnough(costItem) {
	//	this.BuildingLevelUpRes(role.RID(), actorMsg, message.MSG_RESULT_ITEM_NOT_ENOUGH)
	//	return
	//}
	//itemLog := common.NewItemLog(common.BUILDING_LEVEL_UP)
	//role.Building().SetBuildingLevel(req.Id, req.UptoLevel, itemLog)
	//role.SubItemMap(costItem, itemLog)
	//this.BuildingLevelUpRes(role.RID(), actorMsg, message.MSG_RESULT_SUCCESS)
}

func (this *BuildingHandler) MSG_BUILDING_STAR_UP_REQ(actorMsg *actor.ActorMessage) {
	req := actorMsg.Proto().(*message.BuildingStarUpReq)
	role := this.game.PlayerMgr().GetRoleBySession(actorMsg.GateSession())
	building := role.Building().GetBuilding(req.Id)
	if building == nil {
		this.BuildingStarUpRes(role.RID(), actorMsg, message.MSG_RESULT_NO_THIS_BUILDING)
		return
	}
	cfg := config_go.GetBaseBuildingStar(int64(config_global.GetBuildingStarKey(req.Id, req.UptoStar)))
	if cfg == nil { //检测想升到的星级是否存在
		this.BuildingStarUpRes(role.RID(), actorMsg, message.MSG_RESULT_CFG_NO_THIS_PARAM)
		return
	}
	if req.UptoStar <= building.Star {
		this.BuildingStarUpRes(role.RID(), actorMsg, message.MSG_RESULT_CLIENT_WRONG_PARAM)
		return
	}
	var costItem = make(map[int32]int64)
	for star := building.Star; star < req.UptoStar; star++ {
		cfg := config_go.GetBaseBuildingStar(int64(config_global.GetBuildingStarKey(req.Id, star)))
		if cfg == nil {
			this.BuildingStarUpRes(role.RID(), actorMsg, message.MSG_RESULT_CFG_NO_THIS_PARAM)
			return
		}
		for i := 0; i < cfg.Len_STAR_UP_ITEM(); i++ {
			costItem[int32(cfg.Get_STAR_UP_ITEM(i))] += cfg.Get_STAR_UP_NUM(i)
		}
	}

	if !role.ItemsEnough(costItem) {
		this.BuildingStarUpRes(role.RID(), actorMsg, message.MSG_RESULT_ITEM_NOT_ENOUGH)
		return
	}
	itemLog := common.NewItemLog(common.BUILDING_START_UP)
	role.SubItemMap(costItem, itemLog)
	role.Building().SetBuildingStar(req.Id, req.UptoStar, itemLog)
	this.BuildingStarUpRes(role.RID(), actorMsg, message.MSG_RESULT_SUCCESS)
}

func (this *BuildingHandler) MSG_BUILDING_PROP_UP_REQ(actorMsg *actor.ActorMessage) {
	req := actorMsg.Proto().(*message.BuildingPropUpReq)
	role := this.game.PlayerMgr().GetRoleBySession(actorMsg.GateSession())
	building := role.Building().GetBuilding(req.BuildingId)
	if building == nil {
		this.BuildingPropUpRes(role.RID(), actorMsg, message.MSG_RESULT_NO_THIS_BUILDING)
		return
	}

	if building.PropsLevel[req.PropId] == 0 {
		this.BuildingPropUpRes(role.RID(), actorMsg, message.MSG_RESULT_BUILDING_NO_ATTR)
		return
	}

	cfg := config_go.GetBaseBuildingAttr(int64(config_global.GetBuildingPropKey(req.BuildingId, req.PropId, req.UptoLevel)))
	if cfg == nil { //检测想升到的星级是否存在
		this.BuildingPropUpRes(role.RID(), actorMsg, message.MSG_RESULT_CFG_NO_THIS_PARAM)
		return
	}
	if req.UptoLevel <= building.PropsLevel[req.PropId] {
		this.BuildingPropUpRes(role.RID(), actorMsg, message.MSG_RESULT_CLIENT_WRONG_PARAM)
		return
	}
	var costItem = make(map[int32]int64)
	for level := building.PropsLevel[req.PropId]; level < req.UptoLevel; level++ {
		cfg := config_go.GetBaseBuildingAttr(int64(config_global.GetBuildingPropKey(req.BuildingId, req.PropId, level)))
		if cfg == nil {
			this.BuildingPropUpRes(role.RID(), actorMsg, message.MSG_RESULT_CFG_NO_THIS_PARAM)
			return
		}
		for i := 0; i < cfg.Len_LEVEL_UP_ITEM(); i++ {
			costItem[int32(cfg.Get_LEVEL_UP_ITEM(i))] += cfg.Get_LEVEL_UP_NUM(i)
		}
	}

	if !role.ItemsEnough(costItem) {
		this.BuildingPropUpRes(role.RID(), actorMsg, message.MSG_RESULT_ITEM_NOT_ENOUGH)
		return
	}
	itemLog := common.NewItemLog(common.BUILDING_PROP_UP)
	role.SubItemMap(costItem, itemLog)
	role.Building().SetBuildingPropLevel(req.BuildingId, req.PropId, req.UptoLevel, itemLog)
	this.BuildingPropUpRes(role.RID(), actorMsg, message.MSG_RESULT_SUCCESS)
}

func (this *BuildingHandler) BuildingPropUpRes(rid int64, actorMsg *actor.ActorMessage, result message.MSG_RESULT) {
	msg := &message.BuildingPropUpRes{
		Result: result,
	}
	LogResultHandler(rid, result)
	sendInfo := network.NewPbMessage(msg, message.MSG_BUILDING_PROP_UP_RES.Int32())
	system.Send(actorMsg.GateSession(), this.game.GetID(), actorMsg.SourceId(), sendInfo)
}

func (this *BuildingHandler) BuildingLevelUpRes(rid int64, actorMsg *actor.ActorMessage, result message.MSG_RESULT) {
	msg := &message.BuildingLevelUpRes{
		Result: result,
	}
	LogResultHandler(rid, result)
	sendInfo := network.NewPbMessage(msg, message.MSG_BUILDING_LEVEL_UP_RES.Int32())
	system.Send(actorMsg.GateSession(), this.game.GetID(), actorMsg.SourceId(), sendInfo)
}

func (this *BuildingHandler) BuildingStarUpRes(rid int64, actorMsg *actor.ActorMessage, result message.MSG_RESULT) {
	msg := &message.BuildingStarUpRes{
		Result: result,
	}
	LogResultHandler(rid, result)
	sendInfo := network.NewPbMessage(msg, message.MSG_BUILDING_STAR_UP_RES.Int32())
	system.Send(actorMsg.GateSession(), this.game.GetID(), actorMsg.SourceId(), sendInfo)
}
