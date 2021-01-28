package handler

import (
	"root/internal/common"
	"root/internal/config/config_go"
	"root/pkg/actor"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/game/player/slotmachines/slottype"
	"root/servers/internal/message"
)

type SlotHandler struct {
	game interfaces.IGame
}

const FreeCount = 6

func InitSlotHandler(g interfaces.IGame) {
	SlotHandler := &SlotHandler{
		game: g,
	}
	g.RegistMsg(message.MSG_SLOT_BET_REQ.Int32(), SlotHandler.MSG_SLOT_BET_REQ)
	g.RegistMsg(message.MSG_SLOT_SELECT_REQ.Int32(), SlotHandler.MSG_SLOT_SELECT_REQ)
}

// 玩家请求押注
func (this *SlotHandler) MSG_SLOT_BET_REQ(actorMsg *actor.ActorMessage) {
	role := this.roleAssert(actorMsg.GateSession())
	betReq := actorMsg.Proto().(*message.SlotBetReq)
	fields := log.Fields{"RID": role.RID(), "level": role.Lv(), "Glod": role.Gold(), "BetOdds": betReq.GetBetLv()}
	betRes := &message.SlotBetRes{}
	if role.Slot().Select() {
		log.KVs(fields).Warn("player trigger SelectPlay but not select yet!")
		_sendBetResponse(role, message.MSG_RESULT_SLOT_SELECT_TRUE, betRes)
		return
	}
	// 校验玩家金币够不够
	lvConf := config_go.GetSlotCherryBet(betReq.GetBetLv())
	if lvConf == nil {
		log.KVs(fields).Error("can't find lvConf ")
		_sendBetResponse(role, message.MSG_RESULT_FAILED, betRes)
		return
	}
	if role.Gold() < int64(lvConf.BET()) {
		log.KVs(fields).Warn(colorized.Red("player gold not enought"))
		_sendBetResponse(role, message.MSG_RESULT_GOLD_NOT_ENOUGH, betRes)
		return
	}
	// 检查押注金额，是否在配置内
	BetVal := lvConf.BET()

	// 检查等级够不够
	if int64(role.Lv()) < lvConf.LEVEL() {
		fields["tConf.LEVEL()"] = lvConf.LEVEL()
		log.KVs(fields).Error("player level < lvConf.LEVEL()")
		_sendBetResponse(role, message.MSG_RESULT_LEVEL_NOT_ENGOUTH, betRes)
		return
	}

	betRes.BaseData = &message.Base{}
	betRes.BonusData = &message.Bonus{}

	itemLog := common.NewItemLog(common.SLOT_BET)
	role.SubItemMap(map[int32]int64{common.GoldID: BetVal}, itemLog)
	addGold, base, bonus, scatter := role.Slot().Bet(BetVal, int64(role.Lv()), slottype.Trigger_Respin_Base, int(lvConf.POW()), betRes.BaseData, betRes.BonusData)
	role.Slot().SetLastBetLv(betReq.GetBetLv())
	itemmap := map[int32]int64{}
	mapCombine(itemmap, base)
	mapCombine(itemmap, bonus)
	mapCombine(itemmap, scatter)
	log.KVs(fields).KV("bonusItemMap", bonus).KV("scatterItemMap", scatter).KV("baseItemMap", base).Debug("trigger base itemMap")
	role.AddItemMap(itemmap, itemLog)
	role.AddGold(addGold, itemLog)
	betRes.TotalGold = addGold
	betRes.SelectPlay = role.Slot().Select()
	log.KVs(log.Fields{"RID": role.RID(), "betVal": BetVal, "gold": role.Gold(), "addGold": addGold}).
		Debug(colorized.Gray("MSG_SLOT_CHOOSE_BET_REQ"))

	role.Slot().Save()
	_sendBetResponse(role, message.MSG_RESULT_SUCCESS, betRes)

}
func _sendBetResponse(r interfaces.IPlayer, result message.MSG_RESULT, msg *message.SlotBetRes) {
	msg.Result = result
	r.SendMsg(network.NewPbMessage(msg, message.MSG_SLOT_BET_RES.Int32()))
}

// 玩家选择玩法
func (this *SlotHandler) MSG_SLOT_SELECT_REQ(actorMsg *actor.ActorMessage) {
	role := this.roleAssert(actorMsg.GateSession())
	selectReq := actorMsg.Proto().(*message.SlotSelectReq)
	fields := log.Fields{"RID": role.RID(), "level": role.Lv(), "Glod": role.Gold(), "BetOdds": selectReq.GetSelect()}
	selectResp := &message.SlotSelectRes{}
	if !role.Slot().Select() {
		log.KVs(fields).Warn("player can't select")
		_sendSelectResponse(role, message.MSG_RESULT_SLOT_SELECT_FALSE, selectResp)
		return
	}

	selectResp.Select = selectReq.Select
	itemLog := common.NewItemLog(common.SLOT_BET)
	lvConf := config_go.GetSlotCherryBet(role.Slot().LastBetLv())
	if lvConf == nil {
		fields["lastBetLv"] = role.Slot().LastBetLv()
		log.KVs(fields).Error("can't find lvConf ")
		_sendSelectResponse(role, message.MSG_RESULT_FAILED, selectResp)
		return
	}

	if selectReq.Select == 1 { // 选择Bonus玩法
		bonusMsg := &message.Bonus{}
		randrow, Odds := role.Slot().Select1(slottype.Trigger_Respin_Select, bonusMsg)
		addGlod := Odds * int64(lvConf.BET())
		role.AddGold(addGlod, itemLog)
		bonusMsg.BonusGlod = addGlod
		selectResp.RandRow = int64(randrow)
		selectResp.BonusData = bonusMsg
	} else { // Free玩法
		selectResp.Games = []*message.FreeGame{}
		count := FreeCount
		for ; count > 0; count-- {
			baseMsg := &message.Base{}
			bonusMsg := &message.Bonus{}
			addGold, base, bonus, scatter := role.Slot().Bet(lvConf.BET(), int64(role.Lv()), int(lvConf.POW()), slottype.Trigger_Respin_Free, baseMsg, bonusMsg)
			if role.Slot().Select() {
				count += FreeCount
			}
			itemmap := map[int32]int64{}
			mapCombine(itemmap, base)
			mapCombine(itemmap, bonus)
			mapCombine(itemmap, scatter)
			log.KVs(fields).KV("bonusItemMap", bonus).KV("scatterItemMap", scatter).KV("baseItemMap", base).Debug("trigger select itemMap")
			role.AddItemMap(itemmap, itemLog)
			role.AddGold(addGold, itemLog)
			selectResp.Games = append(selectResp.Games, &message.FreeGame{BaseData: baseMsg, BonusData: bonusMsg})
		}
	}
	role.Slot().Save()
	_sendSelectResponse(role, message.MSG_RESULT_SUCCESS, selectResp)
}
func _sendSelectResponse(r interfaces.IPlayer, result message.MSG_RESULT, msg *message.SlotSelectRes) {
	msg.Result = result
	r.SendMsg(network.NewPbMessage(msg, message.MSG_SLOT_SELECT_RES.Int32()))
}

func (this *SlotHandler) roleAssert(gateSession string) interfaces.IPlayer {
	role := this.game.PlayerMgr().GetRoleBySession(gateSession)
	if role == nil {
		log.KV("gateSession", gateSession).ErrorStack(3, "player == nil")
		panic(nil)
	}
	return role
}

// 把目标里的target拷贝到source里
func mapCombine(source, target map[int32]int64) {
	if target == nil || source == nil {
		return
	}
	for k, v := range target {
		if _, ok := source[k]; ok {
			source[k] += v
		} else {
			source[k] = v
		}
	}
}
