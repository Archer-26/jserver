package client

import (
	"root/internal/system"
	"root/pkg/network"
)

type clientSessionHandler struct {
	*Client
	network.BaseNetHandler
}

///////////////////////////////////////// clientSessionHandler /////////////////////////////////////////////
func (this *clientSessionHandler) OnSessionCreated() {
	pointer, _ := this.Load("ext")
	this.Client = pointer.(*Client)
}

func (this *clientSessionHandler) OnSessionClosed() {

}

func (this *clientSessionHandler) OnRecv(bytes []byte) {
	msg := network.NewBytesMessageParse(bytes, this.typemap)
	system.LocalSend("", this.GetID(), "OnRecv", func() {
		system.Send("", "", this.GetID(), msg)
	})
}
