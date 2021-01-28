package client

import (
	"root/pkg/actor"
	"root/pkg/log"
	"root/servers/internal/message"
)

func (this *Client) MSG_S2C_LOGIN_RES(actorMsg *actor.ActorMessage) {
	loginRes := actorMsg.Proto().(*message.LoginRes)
	this.UID = loginRes.UID
	this.RID = loginRes.RID
	this.SID = loginRes.SID
	this.SecurityCode = loginRes.SecurityCode
	this.newPlayer = loginRes.NewPlayer
	this.timestamp = loginRes.Timestamp
	log.KV("S2CLoginRes", loginRes.String()).Info("login success ")
	this.Enter()
}

func (this *Client) MSG_S2C_ENTER_GAME_RES(actorMsg *actor.ActorMessage) {
	enterGameRes := actorMsg.Proto().(*message.EnterGameRes)
	log.KV("S2CEnterGameRes", enterGameRes.String()).Info("enter game")
}

func (this *Client) MSG_NOTIFY_ROLE_INFO(actorMsg *actor.ActorMessage) {
	roleInfo := actorMsg.Proto().(*message.NotifyRoleInfo)
	log.KV("roleInfo", roleInfo.String()).Info("roleInfo ")
}

func (this *Client) MSG_SLOT_BET_RES(actorMsg *actor.ActorMessage) {
	enterGameRes := actorMsg.Proto().(*message.SlotBetRes)
	log.KV("SlotChooseBetRes", enterGameRes.String()).Info("bet resp")
}

func (this *Client) MSG_USE_ITEM_RES(actorMsg *actor.ActorMessage) {
	enterGameRes := actorMsg.Proto().(*message.UseItemRes)
	log.KV("UseItemRes", enterGameRes.String()).Info("useItem resp")
}

func (this *Client) MSG_NOTIFY_ITEMS(actorMsg *actor.ActorMessage) {
	enterGameRes := actorMsg.Proto().(*message.NotifyItems)
	log.KV("NotifyItems", enterGameRes.String()).Info("NotifyItems resp")
}
