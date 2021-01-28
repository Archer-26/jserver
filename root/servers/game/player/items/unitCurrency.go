package items

import (
	"root/internal/common"
	"root/pkg/log"
	"root/servers/game/interfaces"
)

type UnitCurrency struct {
	BaseItem
}

func newUnitCurrency() *UnitCurrency {
	return &UnitCurrency{}
}

func (this *UnitCurrency) Add(role interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog) {
	switch itemId {
	case common.GoldID:
		role.AddGold(itemCount, itemLog)
	case common.Exp:
		role.AddExp(itemCount, itemLog)
	default:
		log.KVs(log.Fields{"rid": role.RID(), "itemId": itemId, "count": itemCount, "itemLog": itemLog}).Error("Add UnitCurrency not implement")
	}
}

func (this *UnitCurrency) Sub(role interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog) {
	switch itemId {
	case common.GoldID:
		role.SubGold(itemCount, itemLog)
	default:
		log.KVs(log.Fields{"rid": role.RID(), "itemId": itemId, "count": itemCount, "itemLog": itemLog}).Error("Sub UnitCurrency not implement")
	}
}

func (this *UnitCurrency) ItemEnough(role interfaces.IPlayer, itemId int32, itemCount int64) bool {
	switch itemId {
	case common.GoldID:
		return role.Gold() >= itemCount
	default:
		log.KVs(log.Fields{"rid": role.RID(), "itemId": itemId, "count": itemCount}).Error("ItemEnoughnot implement")
	}
	return false
}
