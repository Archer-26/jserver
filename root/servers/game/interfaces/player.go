package interfaces

import (
	"root/internal/common"
	"root/pkg/ev"
	"root/pkg/network"
	"root/servers/internal/mysql/i_mysql"
)

type IPlayer interface {
	RID() int64
	UID() int64
	SID() int64
	Name() string
	Lv() int32
	Exp() int64
	Gold() int64
	GateSession() string
	BindSession(gateSession string)
	UnBindSession()
	SyncAllData()
	Online() bool
	SendMsg(msg *network.Message)
	LoginTime() int64
	LogoutTime() int64

	ItemsEnough(itemMap map[int32]int64) bool
	AddItemMap(itemMap map[int32]int64, reason *common.ItemLog)
	SubItemMap(itemMap map[int32]int64, reason *common.ItemLog)

	AddGold(count int64, itemLog *common.ItemLog)
	SubGold(count int64, itemLog *common.ItemLog)
	AddExp(exp int64, itemLog *common.ItemLog)
	Dispatch(ev ev.IEvent)

	// 所有功能对外接口
	Role() IRole
	Storage() IStorage
	Slot() ISlot
	Building() IBuilding
	Island() IIsland
	Dice() IDice
}

const (
	Role = iota //背包
	Storage
	Slot
	Building
	Island
	Dice

	MaxIUnit
)

type IRoleUnit interface {
	Init(rid int64, g IGame)
	Model() i_mysql.IDbTable
	Save()
	Sync() //首次同步所有数据
}
