package network

import (
	"net"
	"root/pkg/abtime"
	"sync/atomic"
)

type SessionType int

const (
	TYPE_TCP SessionType = 1
	TYPE_UDP SessionType = 2
	TYPE_WS  SessionType = 3
)

type INetListener interface {
	IsRunning() bool
	Stop()
}

type INetClient interface {
	NetSession() INetSession
	IsRunning() bool
	Reconnect() bool
	Stop()
}

type INetSession interface {
	Id() uint32
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	RemoteIP() string
	SendMsg([]byte)
	IsRunning() bool
	Stop()
	StoreKV(interface{}, interface{})
	DeleteKV(interface{})
	Load(interface{}) (interface{}, bool)
	CopyStore(INetSession)
	Type() SessionType
}

type INetHandler interface {
	OnSessionCreated()
	OnSessionClosed()
	OnRecv([]byte)
	setSession(INetSession)
}

type BaseNetHandler struct {
	INetSession
}

func (this *BaseNetHandler) setSession(session INetSession) {
	this.INetSession = session
}

var GenNetSessionId = _gen_net_session_id()

func _gen_net_session_id() func() uint32 {
	now := abtime.Now()
	_session_gen_id := uint32(now.Hour()*100000000 + now.Minute()*1000000 + now.Second()*10000)
	return func() uint32 { return atomic.AddUint32(&_session_gen_id, 1) }
}
