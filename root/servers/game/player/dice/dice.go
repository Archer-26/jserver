package dice

import (
	"root/internal/common"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
	"root/servers/internal/models"
	"root/servers/internal/mysql/i_mysql"
)

type (
	Dice struct {
		RID  int64
		game interfaces.IGame
		*models.ModelDice
	}
)

func (this *Dice) Init(rid int64, g interfaces.IGame) {
	this.RID = rid
	this.game = g
	this.ModelDice = models.NewModelDice(rid)
}

func (this *Dice) Sync() {
	DiceInfo := &message.NotifyDiceInfo{
		JsonInfo: this.GameInfo,
	}

	msg := network.NewPbMessage(DiceInfo, message.MSG_NOTIFY_DICE_INFO.Int32())
	this.player().SendMsg(msg)
}

func (this *Dice) Save() {
	this.game.Cache().Save(this.RID, this.ModelDice)
}

// 回存游戏信息
func (this *Dice) SaveGameInfo(msg *message.DiceSaveInfoReq)(res *message.DiceSaveInfoRes) {
	res = &message.DiceSaveInfoRes{Result: message.MSG_RESULT_SUCCESS}

	if !this.player().ItemsEnough(msg.ConsumeItems){
		res.Result = message.MSG_RESULT_ITEM_NOT_ENOUGH
		log.KVs(log.Fields{"rid": this.RID,"items":msg.ConsumeItems}).Debug("Dice ItemsEnough false")
		return
	}

	this.player().SubItemMap(msg.ConsumeItems,common.NewItemLog(common.DICE))
	this.player().AddItemMap(msg.AcquireItems,common.NewItemLog(common.DICE))
	this.GameInfo = msg.JsonInfo
	this.Save()
	log.KVs(log.Fields{"rid": this.RID,"json":msg.JsonInfo,"sub":msg.ConsumeItems,"add":msg.AcquireItems}).Debug("Dice SaveGameInfo")
	return
}

func (this *Dice) Model() i_mysql.IDbTable {
	return this.ModelDice
}

func (this *Dice) player() interfaces.IPlayer {
	return this.game.PlayerMgr().GetRoleByRID(this.RID)
}
