package game

import (
	"root/internal/common"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/iniconfig"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/servers/game/cachemod"
	"root/servers/game/gamedb"
	"root/servers/game/handler"
	"root/servers/game/handler_inner"
	"root/servers/game/interfaces"
	"root/servers/game/player"
)

type (
	Game struct {
		actor.IActor
		playerMgr    *player.Manager
		dbActorID    string
		cache        *cachemod.CacheModels
		msgMap       map[int32]interfaces.RegistFun // HandleMessage没有处理的消息,交给msgMap
		gameStopping bool
	}
)

func (this *Game) Init(a actor.IActor) {
	this.IActor = a

	this.msgMap = make(map[int32]interfaces.RegistFun)
	// game自带的DBActor,处理role相关的所有DB数据
	this.dbActorID = common.DBName(common.Game_Actor, iniconfig.AppId())
	system.Regist(actor.NewActor(this.dbActorID, &gamedb.GameDB{}, actor.SetLocalized()))

	// 初始化业务逻辑需要的属性
	this.playerMgr = player.NewMgr(this)
	this.cache = cachemod.Init(this)

	// 初始化handler
	handler.InitEnterGameHandler(this)
	handler.InitItemHandler(this)
	handler.InitSlotHandler(this)
	handler.InitIslandHandler(this)
	handler.InitDiceHandler(this)
	handler.InitBuildingHandler(this)
	handler_inner.InitHandler(this)
	this.initCmd()
}

// 注册消息回调
func (this *Game) RegistMsg(msgId int32, f interfaces.RegistFun) {
	if _, ok := this.msgMap[msgId]; ok {
		log.KV("msgId", msgId).ErrorStack(3, "regist repeated message")
		return
	}
	this.msgMap[msgId] = f
}

// 返回false表示需要延迟处理，待可以退出时，game自己调用Suspend()
func (this *Game) Stop() bool {
	this.ResetTimerMgr() // 清空计时器
	this.cache.Stop()    // 停止并回存当前所有角色数据
	log.KV("actor", this.GetID()).Info(colorized.Red("Game store data over !!!"))
	return true
}

func (this *Game) HandleMessage(actorMsg *actor.ActorMessage) {
	if this.gameStopping && common.IsGate(actorMsg.SourceId()) {
		return // 如果进入stopping状态，抛弃所有网关消息
	}

	if handle := this.msgMap[actorMsg.MsgId()]; handle != nil {
		handle(actorMsg)
	}
	return
}

///////////////////////////////// 当前game的全局属性 /////////////////////////////////////////////////////
func (this *Game) PlayerMgr() interfaces.IPlayerMgr {
	return this.playerMgr
}
func (this *Game) Cache() interfaces.ICache {
	return this.cache
}

func (this *Game) DBName() string {
	return this.dbActorID
}
func (this *Game) GameStop() bool {
	return this.gameStopping
}
