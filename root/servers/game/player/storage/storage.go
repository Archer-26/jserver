package storage

import (
	"root/internal/common"
	"root/pkg/abtime"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
	"root/servers/internal/models"
	"root/servers/internal/mysql/i_mysql"
)

type (
	Storage struct {
		RID  int64
		game interfaces.IGame
		*models.ModelItems
		modifyItems map[int32]bool
	}
)

func (this *Storage) Init(rid int64, g interfaces.IGame) {
	this.RID = rid
	this.game = g
	this.ModelItems = models.NewModelItems(rid)
	this.modifyItems = make(map[int32]bool)
}

func (this *Storage) AddItem(itemId int32, count int64, itemLog *common.ItemLog) {
	item, ok := this.ItemList[itemId]
	if !ok {
		this.ItemList[itemId] = &models.Item{ItemId: itemId, CreateTime: abtime.Now().Unix()}
		item = this.ItemList[itemId]
	}
	item.ItemCount += count
	log.KVs(log.Fields{"rid": this.RID, "itemId": itemId, "addCount": count, "curCount": item.ItemCount, "itemLog": itemLog}).Info("addItem")
	this.modify(itemId)
}

func (this *Storage) SubItem(itemId int32, count int64, itemLog *common.ItemLog) {
	item, ok := this.ItemList[itemId]
	if !ok {
		log.KVs(log.Fields{"rid": this.RID, "itemId": itemId, "subCount": count, "itemLog": itemLog}).Error("subItem no this ItemId")
		return
	}
	item.ItemCount -= count
	log.KVs(log.Fields{"rid": this.RID, "itemId": itemId, "subCount": count, "curCount": item.ItemCount, "itemLog": itemLog}).Info("subItem")
	if item.ItemCount <= 0 {
		delete(this.ItemList, itemId)
	}
	this.modify(itemId)
}
func (this *Storage) GetItem(itemId int32) int64 {
	item, ok := this.ItemList[itemId]
	if !ok {
		return 0
	}
	return item.ItemCount
}

func (this *Storage) Sync() {
	items := []*message.Item{}
	for _, item := range this.ItemList {
		items = append(items, &message.Item{
			Id:    item.ItemId,
			Count: item.ItemCount,
		})
	}
	allItems := &message.NotifyItems{Items: items}
	msg := network.NewPbMessage(allItems, message.MSG_NOTIFY_ITEMS.Int32())
	this.Role().SendMsg(msg)
	log.KVs(log.Fields{"rid": this.RID, "count": len(items)}).
		Info("sync Storage msg to player")
}

// 回存背包道具数据
func (this *Storage) Save() {
	this.game.Cache().Save(this.RID, this.ModelItems)
}

func (this *Storage) Model() i_mysql.IDbTable {
	return this.ModelItems
}

func (this *Storage) Role() interfaces.IPlayer {
	return this.game.PlayerMgr().GetRoleByRID(this.RID)
}

func (this *Storage) modify(itemId int32) {
	this.modifyItems[itemId] = true
	this.Save()
}

func (this *Storage) NotifyChange() {
	if len(this.modifyItems) <= 0 {
		return
	}
	items := []*message.Item{}
	for itemId, _ := range this.modifyItems {
		msg := &message.Item{Id: itemId}
		if item, ok := this.ItemList[itemId]; ok {
			msg.Count = item.ItemCount
		}
		items = append(items, msg)
	}
	allItems := &message.NotifyItems{Items: items}
	msg := network.NewPbMessage(allItems, message.MSG_NOTIFY_ITEMS.Int32())
	this.Role().SendMsg(msg)
	this.modifyItems = make(map[int32]bool)
}
