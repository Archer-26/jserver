package client

import (
	"math/rand"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/expect"
	"root/pkg/iniconfig"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/internal/message"
	"strings"
	"time"
)

type (
	Client struct {
		conn_addrs []string
		actor.IActor

		SecurityCode string
		RID          int64
		UID          int64
		SID          int64
		passport     string
		uuid         string
		newPlayer    bool
		timestamp    int64
		Gold         int64
		conn         network.INetClient
		typemap      *network.MsgTypeMap
	}
)

func (this *Client) Init(a actor.IActor) {
	this.IActor = a
	this.conn_addrs = strings.Split(iniconfig.String("conn_addr"), ",")
	this.typemap = network.NewMsgTypeMap().InitMsgParser("message", "MSG")
	this.uuid = a.GetID()
	system.RegistCmd(this.GetID(), "login", this.CMD_Login)
	system.RegistCmd(this.GetID(), "enter", this.CMD_Enter)
	system.RegistCmd(this.GetID(), "bet", this.CMD_SlotBet)
	system.RegistCmd(this.GetID(), "s", this.CMD_SlotSelect)
	system.RegistCmd(this.GetID(), "useItem", this.CMD_UseItem)
	system.RegistCmd(this.GetID(), "gm", this.CMD_GM)
	system.RegistCmd(this.GetID(), "blu", this.CMD_BuildingLevelUp)
	system.RegistCmd(this.GetID(), "bsu", this.CMD_BuildingStarUp)
	system.RegistCmd(this.GetID(), "place", this.CMD_Place)
	system.RegistCmd(this.GetID(), "fns", this.CMD_CheckPointFinish)

	system.RegistCmd(this.GetID(), "bpu", this.CMD_BuildingPropUp)

	randGate := rand.Intn(len(this.conn_addrs))
	addr := this.conn_addrs[randGate]
	var ok bool
	this.conn, ok = network.StartTcpClient(addr, &network.StreamCodec{}, &clientSessionHandler{}, network.SetCSMaxRead(1024), network.SetCSExtra(this))
	expect.True(ok, log.Fields{"addr": addr, "msg": "network.StartTcpClient faied"})

	log.KVs(log.Fields{"gate": randGate, "connAddr": this.conn_addrs}).Info("random Gate")
	this.AddTimer(5*time.Second, -1, func(dt int64) {
		ping := network.NewPbMessage(&message.Ping{ClientTimestamp: 999}, message.MSG_PING.Int32())
		this.conn.NetSession().SendMsg(ping.Buffer())
	})

	this.AddTimer(10*time.Second, -1, this.Update)
}

func (this *Client) Login() {
	this.CMD_Login([]string{this.uuid})
}

func (this *Client) Enter() {
	this.CMD_Enter([]string{this.uuid})
}

func (this *Client) Logout() {
	this.Stop()
}

func (this *Client) Stop() bool {
	this.conn.Stop()
	return true
}

func (this *Client) Update(dt int64) {
	//if tools.Probability(10) {
	//	log.InfoOld("stop %v", this.ActorID())
	//	this.Suspend()
	//} else {
	//	this.Login()
	//}
}

func (this *Client) HandleMessage(actorMsg *actor.ActorMessage) {
	log.KVs(log.Fields{"msgId": message.MSG(actorMsg.MsgId()), "info": actorMsg.Proto()}).Info("HandleMessage")
	switch actorMsg.MsgId() {
	case message.MSG_LOGIN_RES.Int32():
		this.MSG_S2C_LOGIN_RES(actorMsg)
	case message.MSG_ENTER_GAME_RES.Int32():
		this.MSG_S2C_ENTER_GAME_RES(actorMsg)
	case message.MSG_NOTIFY_ROLE_INFO.Int32():
		this.MSG_NOTIFY_ROLE_INFO(actorMsg)
	case message.MSG_SLOT_BET_RES.Int32():
		this.MSG_SLOT_BET_RES(actorMsg)
	case message.MSG_USE_ITEM_RES.Int32():
		this.MSG_USE_ITEM_RES(actorMsg)
	case message.MSG_NOTIFY_ITEMS.Int32():
		this.MSG_NOTIFY_ITEMS(actorMsg)
	}
	return
}
