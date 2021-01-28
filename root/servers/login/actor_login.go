package login

import (
	"root/internal/common"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/iniconfig"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/message"
	"root/servers/login/logindb"
)

type (
	Login struct {
		actor.IActor
		dbActorId string
		initiated bool // 初次请求数据
	}
)

func (this *Login) Init(a actor.IActor) {
	this.IActor = a

	this.RegistEvent(actor.AEV_NEW_ACTOR, this.OnNewActor)
	GameMgr.Init(this.IActor)
	InitLoginCmd()

	this.dbActorId = common.DBName(common.Login_Actor, iniconfig.AppId())
	system.Regist(actor.NewActor(this.dbActorId, &logindb.LoginDB{}, actor.SetLocalized()))
	log.Info("watting connect mysql........")
}

// 获取完成所有数据，开启服务
func (this *Login) startService() {
	UserMgr.AllUID(iniconfig.Int32("global_number"))

}

func (this *Login) HandleMessage(actorMsg *actor.ActorMessage) {
	msgId := actorMsg.MsgId()
	switch msgId {
	case inner.INNER_MSG_D2L_ALL_USER_RES.Int32(): // mysql 返回所有user信息
		this.INNER_MSG_D2L_ALL_USER_RES(actorMsg)
	case inner.INNER_MSG_G2L_ROLE_OFFLINE.Int32(): // game 通知玩家离线
		this.INNER_MSG_G2L_ROLE_OFFLINE(actorMsg)
	case message.MSG_LOGIN_REQ.Int32(): // 登录
		this.MSG_TYPE_LOGIN_REQ(actorMsg)
	default:
		return
	}
	return
}

// 事件处理
func (this *Login) OnNewActor(event actor.ActorEvent) {
	switch event.EType() {
	case actor.AEV_NEW_ACTOR:
		evData, ok := event.(actor.Ev_newActor)
		if !ok {
			log.KV("etype", event.EType()).Error("event type error")
			return
		}

		// 成功连接上db后，开始请求所有用户数据
		if this.dbActorId == evData.ActorId && !this.initiated {
			// 向db请求所有用户角色数据
			req := network.NewPbMessage(&inner.L2DAllUserReq{}, inner.INNER_MSG_L2D_ALL_USER_REQ.Int32())
			system.Send("", this.GetID(), this.dbActorId, req)
			this.initiated = true
		} else if common.IsGame(evData.ActorId) {
			GameMgr.SetGameActor(evData.ActorId)
		}
	}
}

func (this *Login) Stop() bool {
	return true
}
