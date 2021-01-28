package logindb

import (
	"fmt"
	"root/internal/common"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/iniconfig"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/mysql"
)

var LoginDBName string

type (
	LoginDB struct {
		actor.IActor
		mysql *mysql.MySQLDatabase
	}
)

func (this *LoginDB) Init(a actor.IActor) {
	this.IActor = a
	LoginDBName = a.GetID()

	this.mysql = &mysql.MySQLDatabase{
		DataBaseName: iniconfig.String(common.LoginDbName),
		Count:        1,
	}
	mysqlName := fmt.Sprintf("%v_%v", a.GetID(), "MySQL")
	sqlActor := actor.NewActor(mysqlName, this.mysql, actor.SetMailBoxSize(5000), actor.SetLocalized())
	this.mysql.IActor = sqlActor // 在注册之前，先把actor对象赋值给处理者，防止actor初始化晚于login的请求消息
	system.Regist(sqlActor)
}

func (this *LoginDB) Stop() bool {
	this.mysql.SuspendStop()
	return true
}

func (this *LoginDB) ActorID() string {
	return this.GetID()
}

func (this *LoginDB) HandleMessage(actorMsg *actor.ActorMessage) {
	switch actorMsg.MsgId() {
	case inner.INNER_MSG_L2D_ALL_USER_REQ.Int32(): // login请求所有用户数据
		this.INNER_MSG_L2D_ALL_USER_REQ(actorMsg)
	case inner.INNER_MSG_L2D_USER_SAVE.Int32():
		this.INNER_MSG_L2D_USER_SAVE(actorMsg)
	default:
		return
	}
	return
}
