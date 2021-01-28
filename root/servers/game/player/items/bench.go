package items

import (
	"root/internal/common"
	"root/internal/config/config_go"
	"root/pkg/log"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
)

const (
	ItemTypeUnitCurrency = 1 //功能数值类
	ItemTypeNormal       = 2 //普通道具
	ItemTypeGift         = 3 //礼包
	ItemTypeBuilding     = 4 //建筑类

)

type IItem interface {
	Add(role interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog)
	Sub(role interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog)
	Use(role interfaces.IPlayer, itemId int32, itemCount int64, extParams *message.UseExtParams, itemLog *common.ItemLog) bool
	ItemEnough(role interfaces.IPlayer, itemId int32, itemCount int64) bool
}

var itemBench = make(map[int32]IItem)

func init() {
	itemBench[ItemTypeUnitCurrency] = newUnitCurrency()
	itemBench[ItemTypeNormal] = newNormal()
	itemBench[ItemTypeGift] = newGift()
	itemBench[ItemTypeBuilding] = newBuilding()

}

func AddItem(unit interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog) {
	baseItem := config_go.GetBaseItem(int64(itemId))
	if baseItem == nil {
		log.KVs(log.Fields{"itemId": itemId}).Error("no this item")
		return
	}
	iItem := itemBench[int32(baseItem.TYPE())]
	if iItem == nil {
		log.KVs(log.Fields{"itemType": baseItem.TYPE()}).Warn("no this itemType")
		iItem = itemBench[ItemTypeNormal]
	}
	iItem.Add(unit, itemId, itemCount, itemLog)
}

func SubItem(unit interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog) {
	baseItem := config_go.GetBaseItem(int64(itemId))
	if baseItem == nil {
		log.KVs(log.Fields{"itemId": itemId}).Error("no this item")
		return
	}
	iItem := itemBench[int32(baseItem.TYPE())]
	if iItem == nil {
		log.KVs(log.Fields{"itemType": baseItem.TYPE()}).Warn("no this itemType")
		iItem = itemBench[ItemTypeNormal]
	}
	iItem.Sub(unit, itemId, itemCount, itemLog)
}

func UseItem(unit interfaces.IPlayer, itemId int32, itemCount int64, extParams *message.UseExtParams) (flag bool) {
	baseItem := config_go.GetBaseItem(int64(itemId))
	if baseItem == nil {
		log.KVs(log.Fields{"itemId": itemId}).Error("no this item")
		return
	}
	iItem := itemBench[int32(baseItem.TYPE())]
	if iItem == nil {
		log.KVs(log.Fields{"itemType": baseItem.TYPE()}).Warn("no this itemType")
		return
	}
	itemLog := common.NewItemLog(common.USE_ITEM)
	if !iItem.Use(unit, itemId, itemCount, extParams, itemLog) {
		return
	}
	iItem.Sub(unit, itemId, itemCount, itemLog)
	return true
}

func ItemEnough(unit interfaces.IPlayer, itemId int32, itemCount int64) (flag bool) {
	baseItem := config_go.GetBaseItem(int64(itemId))
	if baseItem == nil {
		log.KVs(log.Fields{"itemId": itemId}).Error("no this item")
		return
	}
	iItem := itemBench[int32(baseItem.TYPE())]
	if iItem == nil {
		log.KVs(log.Fields{"itemType": baseItem.TYPE()}).Warn("no this itemType")
		return
	}
	if !iItem.ItemEnough(unit, itemId, itemCount) {
		return
	}
	return true
}
