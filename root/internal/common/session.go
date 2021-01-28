package common

import (
	"fmt"
	"root/pkg/log"
	"strconv"
	"strings"
)

func GateSession(actorID string, sessionId uint32) string {
	return fmt.Sprintf("%v:%v", actorID, sessionId)
}

func SplitGateSession(gateSession string) (actorId string, sessionId uint32) {
	strs := strings.Split(gateSession, ":")
	if len(strs) != 2 {
		log.KV("gateSession", gateSession).ErrorStack(3, "SplitGateSession error")
		panic(nil)
	}
	actorId = strs[0]
	sint, e := strconv.Atoi(strs[1])
	if e != nil {
		log.KV("gateSession", gateSession).ErrorStack(3, "SplitGateSession error")
		panic(nil)
	}
	sessionId = uint32(sint)
	return
}
