package actor

import (
	"fmt"
	"root/pkg/log"
	"sync"
)

/*
	actor 内部抛出的事件，业务层有需要就监听回调
*/

const (
	AEV_NEW_ACTOR         ActorEventType = iota //	新的actor
	AEV_REMOTE_INEXISTENT                       //	给未连接到的远端actor发消息
)

type (
	ActorEventType uint32
	ActorEvent     interface {
		EType() ActorEventType
	}
	Callback     func(ActorEvent)
	EvDispatcher struct {
		sys *ActorSystem
		sync.RWMutex
		ev2listeners map[ActorEventType]map[string]Callback
	}

	Ev_newActor struct {
		ActorId string
		Remote  int8
	}

	Ev_remoteInexistent struct {
		ActorId string
	}
)

func (this Ev_newActor) EType() ActorEventType         { return AEV_NEW_ACTOR }
func (this Ev_remoteInexistent) EType() ActorEventType { return AEV_REMOTE_INEXISTENT }

func NewActorEvents(as *ActorSystem) *EvDispatcher {
	return &EvDispatcher{sys: as, ev2listeners: make(map[ActorEventType]map[string]Callback)}
}

// 注册actor事件
func (s *EvDispatcher) registEvent(actorId string, ev ActorEventType, f Callback) {
	if f == nil {
		log.KV("actorId", actorId).KV("ev", ev).WarnStack(1, "callback is nil")
		return
	}

	s.Lock()
	defer s.Unlock()

	if s.ev2listeners[ev] == nil {
		s.ev2listeners[ev] = make(map[string]Callback)
	}
	s.ev2listeners[ev][actorId] = f
}

// 取消actor事件
func (s *EvDispatcher) cancelEvent(actorId string, et ActorEventType) {
	s.Lock()
	defer s.Unlock()

	if s.ev2listeners[et] != nil {
		delete(s.ev2listeners[et], actorId)
	}
}

// 事件触发
func (s *EvDispatcher) _dispatchEvent(event ActorEvent) {
	s.RLock()
	defer s.RUnlock()

	if listeners := s.ev2listeners[event.EType()]; listeners != nil {
		for actorId, fn := range listeners {
			info := fmt.Sprintf("[%v] process event:[%v]", actorId, event.EType())
			s.sys.LocalSend("", actorId, info, func() { fn(event) })
		}
	}
}
