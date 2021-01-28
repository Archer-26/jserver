package items

import (
	"root/internal/common"
	"root/pkg/log"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
)

//背包中的普通道具
type Normal struct {
}

func newNormal() *Normal {
	return &Normal{}
}

func (this *Normal) Add(role interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog) {
	role.Storage().AddItem(itemId, itemCount, itemLog)
}

func (this *Normal) Sub(role interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog) {
	role.Storage().SubItem(itemId, itemCount, itemLog)
}

func (this *Normal) Use(role interfaces.IPlayer, itemId int32, itemCount int64, extParams *message.UseExtParams, itemLog *common.ItemLog) bool {
	log.KVs(log.Fields{"rid": role.RID(), "itemId": itemId, "count": itemCount}).Error("item Can't use")
	return false
}

func (this *Normal) ItemEnough(role interfaces.IPlayer, itemId int32, itemCount int64) bool {
	return role.Storage().GetItem(itemId) >= itemCount
}
