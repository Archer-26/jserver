package player

import (
	"github.com/vmihailenco/msgpack"
	"root/internal/common"
	"root/internal/system"
	"root/pkg/abtime"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/internal/inner_message/inner"
	"time"
)

type (
	Manager struct {
		game            interfaces.IGame
		playerByRID     map[int64]interfaces.IPlayer
		playerBySession map[string]interfaces.IPlayer
		playerByUID     map[int64]interfaces.IPlayer

		callbackId  int64
		callbackFun map[int64]callbackInfo
	}

	callbackInfo struct {
		fun  func(interfaces.WholeRoleBytes) // 完整的角色信息，包含所需所有功能数据
		time int64                           // 计算超时
		rid  int64
	}
)

func NewMgr(game interfaces.IGame) *Manager {
	mgr := &Manager{
		game:            game,
		playerByUID:     make(map[int64]interfaces.IPlayer),
		playerByRID:     make(map[int64]interfaces.IPlayer),
		playerBySession: make(map[string]interfaces.IPlayer),
		callbackFun:     make(map[int64]callbackInfo),
		callbackId:      abtime.Milliseconds(),
	}
	game.AddTimer(10*time.Minute, -1, mgr.TenMinUpdate)
	return mgr
}

// 玩家第一次登录创建新角色
func (this *Manager) NewRole(gateSession string, RID, UID, SID int64, Passport, PassportType string) interfaces.IPlayer {
	newPlayer := newPLayer(this.game, gateSession)
	// 初始化角色的所有功能模块
	for _, v := range newPlayer.roleUnits {
		v.Init(RID, this.game)
	}
	newPlayer.assignUnits()

	// todo 给新玩家的一些初始化数据，在这里添加
	newPlayer.Role().BaseInfo(UID, SID, Passport, PassportType, "")

	newPlayer.saveAll()
	this.RelateSession2Role(gateSession, newPlayer)
	return newPlayer
}

// 通过数据库加载mod构建角色
func (this *Manager) RoleMod(gateSession string, RID int64, mod interfaces.WholeRoleBytes) interfaces.IPlayer {
	role := newPLayer(this.game, gateSession)
	for _, v := range role.roleUnits {
		v.Init(RID, this.game)
		if bytes, ok := mod[v.Model().ModelName()]; ok {
			_checkErr(msgpack.Unmarshal(bytes, v.Model()), RID)
		} else {
			//老玩家新加的功能会导致db没有参数
			log.KV("mdName", v.Model().ModelName()).KV("RID", RID).
				Warn("md was not found bytes,inspect database table")
		}
	}

	role.assignUnits()
	this.RelateSession2Role(gateSession, role)
	return role
}

func (this *Manager) RelateSession2Role(gateSession string, role interfaces.IPlayer) {
	// 如果角色存在，并且之前绑定过session，顶号处理
	if erole := this.playerByRID[role.RID()]; erole != nil && erole.GateSession() != "" {
		delete(this.playerBySession, erole.GateSession())
		log.KVs(log.Fields{"old gateSession": erole.GateSession(), "new gateSession": gateSession, "RID": role.RID(), "UID": role.UID()}).
			Warn("new session instead of old gateSession")
	}
	this.playerByUID[role.UID()] = role
	this.playerByRID[role.RID()] = role
	this.playerBySession[gateSession] = role
}

func (this *Manager) GetRoleByRID(RID int64) interfaces.IPlayer {
	return this.playerByRID[RID]
}
func (this *Manager) GetRoleBySession(gateSession string) interfaces.IPlayer {
	return this.playerBySession[gateSession]
}
func (this *Manager) GetRoleByUID(UID int64) interfaces.IPlayer {
	return this.playerByUID[UID]
}

func (this *Manager) SetOfflineRole(gateSession string) {
	role := this.playerBySession[gateSession]
	if role == nil {
		log.KVs(log.Fields{"gateSession": gateSession}).Warn("re-login caused by situation 1")
		return // 正常情况,旧的gateSession可能已经被RelateSession2Role()顶号删除
	}
	// 通知login玩家离线
	msg := network.NewPbMessage(&inner.G2LRoleOffline{UID: role.UID(), RID: role.RID()}, inner.INNER_MSG_G2L_ROLE_OFFLINE.Int32())
	system.Send(gateSession, this.game.GetID(), common.Login_Actor, msg)

	role.UnBindSession()
	delete(this.playerBySession, gateSession)

	log.KVs(log.Fields{"RID": role.RID(), "UID": role.UID(), "gateSession": gateSession}).
		Info("SetOfflineRole")
}

func (this *Manager) TenMinUpdate(dt int64) {
	now := abtime.Milliseconds()

	// GetRoleByAsync超时处理
	for callbackid, info := range this.callbackFun {
		if now-info.time > int64(time.Minute*5) {
			log.KVs(log.Fields{"rid": info.rid, "callbackid": callbackid}).
				Warn("GetRoleByAsync timeout")
			delete(this.callbackFun, callbackid)
		}
	}
}

func (this *Manager) GetRoleByAsync(RID int64, callback func(interfaces.WholeRoleBytes)) {
	this.callbackId++
	msgReq := network.NewPbMessage(&inner.G2DRoleReq{CallbackId: this.callbackId, RID: RID}, inner.INNER_MSG_G2D_ROLE_REQ.Int32())
	system.Send("", this.game.GetID(), this.game.DBName(), msgReq)

	this.callbackFun[this.callbackId] = callbackInfo{
		fun:  callback,
		time: abtime.Milliseconds(),
	}
}

// 收到role数据，回调callback
func (this *Manager) CallbackFun(callbackId, RID int64, mod interfaces.WholeRoleBytes) {
	fun, ok := this.callbackFun[callbackId]
	if !ok {
		log.KV("callbackId", callbackId).Warn("can not find callback fun ")
		return
	}
	if mod == nil {
		log.KV("callbackId", callbackId).Info("GetRoleByAsync mod is nil")
	}
	fun.fun(mod)
}

func _checkErr(e error, RID int64) {
	if e != nil {
		log.KVs(log.Fields{"RID": RID, "err": e}).ErrorStack(3, "_checkErr")
		return
	}
}
