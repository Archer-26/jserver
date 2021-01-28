package common

import (
	"fmt"
	"regexp"
	"root/pkg/log"
)

// 静态Actor类型,全局唯一
const (
	//
	Login_Actor = "Login"
)

// 动态增删的actor,会有多个 game1、game2
const (
	Game_Actor = "Game"
	Gate_Actor = "Gate"
	DB_Actor   = "DB"
)

func GameName(id int32) string {
	return fmt.Sprintf("%v%v_Actor", Game_Actor, id)
}

// 匹配game 按照固定格式匹配
func IsGame(actorId string) bool {
	match, e := regexp.MatchString("Game([0-9]+)_Actor", actorId)
	if e != nil {
		log.KV("actorId", actorId).ErrorStack(3, "error")
		return false
	}
	return match
}

func DBName(prex string, id int32) string {
	return fmt.Sprintf("%v%v%v_Actor", prex, DB_Actor, id)
}

func GateName(id int32) string {
	return fmt.Sprintf("%v%v_Actor", Gate_Actor, id)
}

func IsGate(actorId string) bool {
	match, e := regexp.MatchString("Gate([0-9]+)_Actor", actorId)
	if e != nil {
		log.KV("actorId", actorId).ErrorStack(3, "error")
		return false
	}
	return match
}
