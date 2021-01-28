package gamedb

import (
	"fmt"
	"root/internal/common"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/iniconfig"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/mysql"
)

const databaseCount = 2

type (
	GameDB struct {
		actor.IActor
		mysql *mysql.MySQLDatabase
	}
)

func (this *GameDB) Init(a actor.IActor) {
	this.IActor = a
	this.mysql = &mysql.MySQLDatabase{
		DataBaseName: iniconfig.String(common.GameDbName),
		Count:        databaseCount,
	}
	mysqlName := fmt.Sprintf("%v_%v", a.GetID(), "MySQL")
	system.Regist(actor.NewActor(mysqlName, this.mysql, actor.SetMailBoxSize(5000), actor.SetLocalized()))
}

func (this *GameDB) Stop() bool {
	return false
}

func (this *GameDB) ActorID() string {
	return this.GetID()
}

func (this *GameDB) HandleMessage(actorMsg *actor.ActorMessage) {
	switch actorMsg.MsgId() {
	case inner.INNER_MSG_G2D_ROLE_REQ.Int32(): // game请求角色数据
		this.INNER_MSG_G2D_ROLE_REQ(actorMsg)
	case inner.INNER_MSG_G2D_MODEL_SAVE.Int32(): // 存储数据
		this.INNER_MSG_ALL2D_MODEL_SAVE(actorMsg)
	case inner.INNER_MSG_G2D_GAME_STOP.Int32(): // game关闭信号
		this.INNER_MSG_G2D_GAME_STOP(actorMsg)
	default:
		return
	}
	return
}
