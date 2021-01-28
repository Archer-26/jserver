package player

import (
	"root/internal/common"
	"root/internal/system"
	"root/pkg/abtime"
	"root/pkg/ev"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/game/player/building"
	"root/servers/game/player/dice"
	"root/servers/game/player/island"
	"root/servers/game/player/items"
	"root/servers/game/player/role"
	"root/servers/game/player/slotmachines"
	"root/servers/game/player/storage"
	"root/servers/internal/message"
	"strings"
)

type (
	Player struct {
		game        interfaces.IGame
		roleUnits   [interfaces.MaxIUnit]interfaces.IRoleUnit
		role        *role.Role
		slot        *slotmachines.Slot
		storage     *storage.Storage
		building    *building.Building
		island      *island.Island
		dice      	*dice.Dice
		gateSession string
		loginTime   int64
		logoutTime  int64

		events ev.IEventDispatcher
	}
)
func newPLayer(g interfaces.IGame, gateSession string) *Player {
	role := &Player{
		game:        g,
		gateSession: gateSession,
		events:      ev.New(),
	}
	role.prepareUnits()

	return role
}

func (this *Player) prepareUnits() {
	this.roleUnits[interfaces.Role] = &role.Role{}
	this.roleUnits[interfaces.Storage] = &storage.Storage{}
	this.roleUnits[interfaces.Slot] = &slotmachines.Slot{}
	this.roleUnits[interfaces.Building] = &building.Building{}
	this.roleUnits[interfaces.Island] = &island.Island{}
	this.roleUnits[interfaces.Dice] = &dice.Dice{}

	// todo .............
}

func (this *Player) assignUnits() {
	this.role = this.roleUnits[interfaces.Role].(*role.Role)
	this.storage = this.roleUnits[interfaces.Storage].(*storage.Storage)
	this.slot = this.roleUnits[interfaces.Slot].(*slotmachines.Slot)
	this.building = this.roleUnits[interfaces.Building].(*building.Building)
	this.island = this.roleUnits[interfaces.Island].(*island.Island)
	this.dice = this.roleUnits[interfaces.Dice].(*dice.Dice)
	// todo .............
}

func (this *Player) UID() int64          { return this.role.UID }
func (this *Player) RID() int64          { return this.role.RID }
func (this *Player) SID() int64          { return this.role.SID }
func (this *Player) Name() string        { return this.role.Name }
func (this *Player) Lv() int32           { return this.role.Level }
func (this *Player) Exp() int64          { return this.role.Exp }
func (this *Player) Gold() int64         { return this.role.Gold }
func (this *Player) GateSession() string { return this.gateSession }
func (this *Player) Online() bool        { return this.loginTime > this.logoutTime }
func (this *Player) LoginTime() int64    { return this.loginTime }
func (this *Player) LogoutTime() int64   { return this.logoutTime }
func (this *Player) SendMsg(msg *network.Message) {
	if this.GateSession() == "" {
		return
	}
	gate, _ := common.SplitGateSession(this.GateSession())
	system.Send(this.GateSession(), this.game.GetID(), gate, msg)
}

func (this *Player) BindSession(gateSession string) {
	if gateSession == "" {
		log.KVs(log.Fields{"new gateSession": gateSession, "old gateSession": this.gateSession, "RID": this.RID(), "UID": this.UID()}).
			Error("BindSession error")
		return
	}
	// 格式校验
	if len(strings.Split(gateSession, ":")) != 2 {
		log.KVs(log.Fields{"gateSession": this.gateSession, "RID": this.RID(), "UID": this.UID()}).
			ErrorStack(3, "md gateSession split error")
		return
	}
	this.gateSession = gateSession
	this.loginTime = abtime.Milliseconds()
}

func (this *Player) UnBindSession() {
	this.logoutTime = abtime.Milliseconds()
	this.gateSession = ""
}

func (this *Player) saveAll() {
	for _, iUnit := range this.roleUnits {
		iUnit.Save()
	}
}

// 给玩家同步所有数据
func (this *Player) SyncAllData() {
	for _, iUnit := range this.roleUnits {
		iUnit.Sync()
	}
}

func (this *Player) ToRoleInfo() *message.NotifyRoleInfo {
	return &message.NotifyRoleInfo{
		Name:  this.Name(),
		IconW: "",
		Level: this.Lv(),
		Exp:   this.Exp(),
		Gold:  this.Gold(),
	}
}

func (this *Player) Role() interfaces.IRole {
	return this.role
}

func (this *Player) Slot() interfaces.ISlot {
	return this.slot
}

func (this *Player) Storage() interfaces.IStorage {
	return this.storage
}

func (this *Player) Building() interfaces.IBuilding {
	return this.building
}
func (this *Player) Island() interfaces.IIsland {
	return this.island
}

func (this *Player) Dice() interfaces.IDice {
	return this.dice
}

func (this *Player) addItem(itemId int32, itemCount int64, reason *common.ItemLog) {
	items.AddItem(this, itemId, itemCount, reason)
}

func (this *Player) subItem(itemId int32, itemCount int64, itemLog *common.ItemLog) {
	items.SubItem(this, itemId, itemCount, itemLog)
}

func (this *Player) AddItemMap(itemMap map[int32]int64, reason *common.ItemLog) {
	if len(itemMap) <= 0 {
		return
	}
	for itemId, count := range itemMap {
		this.addItem(itemId, count, reason)
	}
	this.Storage().NotifyChange()
}

func (this *Player) ItemsEnough(itemMap map[int32]int64) bool {
	if len(itemMap) <= 0 {
		return false
	}
	for itemId, count := range itemMap {
		if !items.ItemEnough(this, itemId, count) {
			return false
		}
	}
	return true
}

func (this *Player) SubItemMap(itemMap map[int32]int64, reason *common.ItemLog) {
	if len(itemMap) <= 0 {
		return
	}
	for itemId, count := range itemMap {
		this.subItem(itemId, count, reason)
	}
	this.Storage().NotifyChange()
}

func (this *Player) AddGold(count int64, itemLog *common.ItemLog) {
	this.role.AddGold(count, itemLog)
}

func (this *Player) SubGold(count int64, itemLog *common.ItemLog) {
	this.role.SubGold(count, itemLog)
}

func (this *Player) AddExp(exp int64, itemLog *common.ItemLog) {
	this.role.AddExp(exp, itemLog)
}

func (this *Player) Dispatch(ev ev.IEvent) {
	this.events.Dispatch(ev)
}
