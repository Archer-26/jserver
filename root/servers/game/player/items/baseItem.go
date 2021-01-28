package items

import (
	"root/internal/common"
	"root/pkg/log"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
)

//基础Item 返回所有错误操作
type BaseItem struct {
}

func (this *BaseItem) Add(role interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog) {
	log.KVs(log.Fields{"rid": role.RID(), "itemId": itemId, "count": itemCount}).Error("item Can't Add")
}

func (this *BaseItem) Sub(role interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog) {
	log.KVs(log.Fields{"rid": role.RID(), "itemId": itemId, "count": itemCount}).Error("item Can't Sub")
}

func (this *BaseItem) Use(role interfaces.IPlayer, itemId int32, itemCount int64, extParams *message.UseExtParams, itemLog *common.ItemLog) bool {
	log.KVs(log.Fields{"rid": role.RID(), "itemId": itemId, "count": itemCount}).Error("item Can't Use")
	return false
}

func (this *BaseItem) ItemEnough(role interfaces.IPlayer, itemId int32, itemCount int64) bool {
	log.KVs(log.Fields{"rid": role.RID(), "itemId": itemId, "count": itemCount}).Error("item Can't ItemEnough")
	return false
}
