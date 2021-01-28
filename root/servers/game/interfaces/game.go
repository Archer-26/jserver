package interfaces

import "root/pkg/actor"

type IGame interface {
	actor.IActor
	PlayerMgr() IPlayerMgr
	Cache() ICache
	RegistMsg(msgId int32, f RegistFun)
	DBName() string
}

// 消息注册,消息优先流给HandleMessage, HandleMessage不处理的,再交给这里
type RegistFun func(*actor.ActorMessage)
