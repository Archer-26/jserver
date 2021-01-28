package interfaces

import "root/internal/common"

type IStorage interface {
	AddItem(itemId int32, count int64, itemLog *common.ItemLog)
	SubItem(itemId int32, count int64, itemLog *common.ItemLog)
	GetItem(itemId int32) int64
	NotifyChange()
}
