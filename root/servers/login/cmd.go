package login

import (
	"encoding/json"
	"fmt"
	"root/internal/common"
	"root/internal/system"
	"root/pkg/expect"
	"root/pkg/log"
	"strconv"
)

func InitLoginCmd() {
	system.RegistCmd(common.Login_Actor, "game", avialableGame)
	system.RegistCmd(common.Login_Actor, "usercount", userCount)
	system.RegistCmd(common.Login_Actor, "user", searchUser)
}

func avialableGame(param []string) {
	fmt.Printf("可分配的game:\n")
	for k, v := range GameMgr.gameActors {
		info := fmt.Sprintf(" %v 人数:%v\n", k, v.Player_number)
		log.Info(info)
	}
}

func userCount(param []string) {
	info := fmt.Sprintf("当前所有用户数量:%v\n", UserMgr.UserLen())
	log.Info(info)
}

func searchUser(param []string) {
	expect.True(len(param) == 1)
	i, e := strconv.Atoi(param[0])
	expect.Nil(e,log.Fields{"param[0]":param[0]})

	uid := int64(i)
	user := UserMgr.User(uid)
	v, _ := json.Marshal(user)
	info := fmt.Sprintf("user:%v\n", string(v))
	log.Info(info)
}
