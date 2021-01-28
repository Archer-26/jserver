package webconsole

import (
	"root/pkg/actor"
)

type (
	WebConsole struct {
		actor.IActor
	}
)

func (this *WebConsole) Init(a actor.IActor) {
	this.IActor = a
}

func (this *WebConsole) Stop() bool {

	return true
}

func (this *WebConsole) HandleMessage(actorMsg *actor.ActorMessage) {
	return
}
