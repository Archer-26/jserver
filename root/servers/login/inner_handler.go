package login

import (
	"github.com/vmihailenco/msgpack"
	"root/pkg/actor"
	"root/pkg/expect"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/models"
)

// mysql 返回所有用户信息
func (this *Login) INNER_MSG_D2L_ALL_USER_RES(actorMsg *actor.ActorMessage) {
	res := actorMsg.Proto().(*inner.D2LAllUserRes)
	if res.Data == nil {
		log.Info(colorized.Cyan("all user receive finish!"))
		this.startService()
		return
	}

	allUser := []*models.ModelUser{}
	err := msgpack.Unmarshal(res.Data, &allUser)
	expect.Nil(err)

	for _, v := range allUser {
		UserMgr.users[v.UID] = v
		for rid, info := range v.Roles {
			UserMgr.roles[rid] = info
		}
		for passwd, passport := range v.Passports {
			UserMgr.passports[passwd] = passport
		}
	}
	log.KV("len", len(allUser)).Info(colorized.Cyan("get all user's info"))
}

// game 通知玩家离线
func (this *Login) INNER_MSG_G2L_ROLE_OFFLINE(actorMsg *actor.ActorMessage) {
	res := actorMsg.Proto().(*inner.G2LRoleOffline)
	user := UserMgr.users[res.UID]
	expect.True(user != nil,log.Fields{"UID":res.UID})

	if user.GateSession == actorMsg.GateSession() {
		user.GateSession = "" // 正常离线
	} else {
		// 出现这种情况，说明 玩家新的session登录成功后 < 旧session断开 < 新session完成enterGame消息
		log.KVs(log.Fields{"UID": res.UID, "gateSession": actorMsg.GateSession()}).
			Warn("re-login caused by situation 2")
	}

	log.KV("UID", res.UID).Info("game inform user offline")
}
