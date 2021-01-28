package system

import (
	"root/internal/coder"
	"root/pkg/actor"
	"root/pkg/iniconfig"
	"root/pkg/network"
)

var global *actor.ActorSystem

func InitActorSystem() {
	actor.InitCmd()
	global = actor.NewActorSystem(
		actor.Coder(coder.NewMessageCoder()),
		actor.ActorAddr(iniconfig.String("actor_port")),
		actor.EtcdAddr(iniconfig.String("etcd_addr")),
		actor.EtcdPrefix(iniconfig.String("etcd_prefix")),
	)
}

func Regist(a actor.IActor) {
	global.Regist(a)
}

func GateSpecificRemoteCoder() {
	global.RemoteCoder(coder.NewBytesCoder())
}

func Startup() {
	global.Startup()
}

func RegistCmd(actorId, cmd string, f func([]string)) {
	global.RegistCmd(actorId, cmd, f)
}

func LocalSend(sourceId, targetId, funName string, f func()) {
	global.LocalSend(sourceId, targetId, funName, f)
}

func Send(gateSession, sourceId, targetId string, pbMsg *network.Message) {
	global.Send(gateSession, sourceId, targetId, pbMsg)
}

func LocalActorFuzzy(f func(reg string) bool) []actor.IActor {
	return global.LocalActorFuzzy(f)
}
