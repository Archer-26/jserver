package ev

import (
	"root/pkg/log"
	"root/pkg/tools"
)

/*
	单线程事件机制
*/
type (
	EventType        uint32
	IEvent           interface{ EType() EventType }
	IEventListener   interface{ OnEvent(IEvent) }
	IEventDispatcher interface {
		Dispatch(event IEvent)
		AddEventListener(eventType EventType, listener IEventListener)
		RemoveListenerByEType(eventType EventType, listener IEventListener)
		RemoveListenerAll(listener IEventListener)
	}

	eventDispatcher struct {
		listeners map[EventType][]IEventListener
	}
)

func New() IEventDispatcher {
	return &eventDispatcher{
		listeners: make(map[EventType][]IEventListener),
	}
}

func (this *eventDispatcher) Dispatch(event IEvent) {
	list, ok := this.listeners[event.EType()]
	if !ok {
		//log.KV("etype", event.EType()).WarnStack(1, "event has no listener")
		return
	}

	for i := range list {
		tools.Try(func() { list[i].OnEvent(event) }, nil)
	}
}

func (this *eventDispatcher) AddEventListener(eventType EventType, listener IEventListener) {
	_, ok := this.listeners[eventType]
	if !ok {
		this.listeners[eventType] = make([]IEventListener, 0)
	}

	list := this.listeners[eventType]
	for i := range list {
		if list[i] == listener {
			log.KV("etype", eventType).Error("AddEventListener repeated")
			return
		}
	}
	this.listeners[eventType] = append(list, listener)
}

func (this *eventDispatcher) RemoveListenerByEType(eventType EventType, listener IEventListener) {
	list, ok := this.listeners[eventType]
	if !ok {
		log.KV("eventType", eventType).Error("RemoveListenerByEType failed")
		return
	}
	for i := range list {
		if list[i] == listener {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	this.listeners[eventType] = list
}

func (this *eventDispatcher) RemoveListenerAll(listener IEventListener) {
	for typ := range this.listeners {
		this.RemoveListenerByEType(typ, listener)
	}
}
