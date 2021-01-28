package gate

import (
	"root/internal/common"
	"root/internal/system"
	"root/pkg/abtime"
	"root/pkg/expect"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/message"
)

type UserSessionHandler struct {
	GateHandler *Gate
	GameActor   string // 处理当前session的game
	Alive       bool
	network.BaseNetHandler
}

func (this *UserSessionHandler) OnSessionCreated() {
	pointer, ok := this.Load("ext")
	expect.True(ok)

	this.Alive = true
	this.GateHandler = pointer.(*Gate)
	// 这里只做session映射，等待客户端请求登录
	system.LocalSend("", this.GateHandler.GetID(), "OnSessionCreated", func() {
		this.GateHandler.sessions[this.Id()] = this
	})
}

func (this *UserSessionHandler) OnSessionClosed() {
	if this.GameActor != "" {
		// 连接断开，通知game
		gateSession := common.GateSession(this.GateHandler.GetID(), this.INetSession.Id())
		msg := network.NewPbMessage(&inner.GT2GSessionClosed{GateSession: gateSession}, inner.INNER_MSG_GT2G_SESSION_CLOSED.Int32())
		system.Send(gateSession, this.GateHandler.GetID(), this.GameActor, msg)
	}

	system.LocalSend("", this.GateHandler.GetID(), "OnSessionClosed", func() {
		delete(this.GateHandler.sessions, this.Id())
	})
}

func (this *UserSessionHandler) OnRecv(data []byte) {
	expect.True(len(data) >= 4,log.Fields{"len(data)": len(data), "session": this.INetSession.Id()})

	msgId := network.Byte4ToUint32(data[:4])
	// gate特殊处理用户心跳
	if msgId == message.MSG_PING.UInt32() {
		ping := network.NewBytesMessageParse(data, this.GateHandler.typemap).Proto().(*message.Ping)
		pong := network.NewPbMessage(&message.Pong{
			ClientTimestamp: ping.ClientTimestamp,
			ServerTimestamp: abtime.Milliseconds(),
		}, message.MSG_PONG.Int32())
		this.INetSession.SendMsg(pong.Buffer())
		this.Alive = true
		return
	}

	log.KV("MsgId", msgId).Info("UserSessionHandler OnRecv msgID")
	var msg *network.Message
	if this.GateHandler.Parse {
		msg = network.NewBytesMessageParse(data, this.GateHandler.typemap)
	} else {
		msg = network.NewBytesMessage(data)
	}

	gateSession := common.GateSession(this.GateHandler.GetID(), this.INetSession.Id())
	if message.MSG_LOGIN_SEGMENT_BEGIN.UInt32() <= msgId && msgId <= message.MSG_LOGIN_SEGMENT_END.UInt32() {
		system.Send(gateSession, this.GateHandler.GetID(), common.Login_Actor, msg)
	} else if message.MSG_GAME_SEGMENT_BEGIN.UInt32() <= msgId && msgId <= message.MSG_GAME_SEGMENT_END.UInt32() {
		expect.True(this.GameActor != "",log.Fields{"session": this.Id(), "msgId": msgId})
		system.Send(gateSession, this.GateHandler.GetID(), this.GameActor, msg)
	}
}
