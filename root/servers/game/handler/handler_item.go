package handler

import (
	"flag"
	"root/internal/common"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/game/player/items"
	"root/servers/internal/message"
	"runtime/debug"
	"strconv"
	"strings"
)

type StorageHandler struct {
	game interfaces.IGame
}

func InitItemHandler(g interfaces.IGame) {
	storageHandler := &StorageHandler{
		game: g,
	}
	g.RegistMsg(message.MSG_USE_ITEM_REQ.Int32(), storageHandler.MSG_USE_ITEM_REQ)
	g.RegistMsg(message.MSG_SELL_ITEM_REQ.Int32(), storageHandler.MSG_SELL_ITEM_REQ)
	g.RegistMsg(message.MSG_GM_COMMAND_REQ.Int32(), storageHandler.MSG_GM_COMMAND_REQ)
}

var gmCommands = make(map[string]func(role interfaces.IPlayer, params []string) bool)

func init() {
	gmCommands["addItem"] = addItem //addItem itemId itemCount //增加某种类型的道具   Arg1:道具类型 Arg2:数量
}

func (this *StorageHandler) MSG_USE_ITEM_REQ(actorMsg *actor.ActorMessage) {
	useItemReq := actorMsg.Proto().(*message.UseItemReq)
	role := this.game.PlayerMgr().GetRoleBySession(actorMsg.GateSession())

	if !role.ItemsEnough(map[int32]int64{useItemReq.ItemId: useItemReq.Count}) {
		this.UseItemRes(role.RID(), actorMsg, message.MSG_RESULT_ITEM_NOT_ENOUGH)
		return
	}
	result := message.MSG_RESULT_FAILED
	if items.UseItem(role, useItemReq.ItemId, useItemReq.Count, useItemReq.Params) {
		result = message.MSG_RESULT_SUCCESS
	}
	this.UseItemRes(role.RID(), actorMsg, result)
}

func (this *StorageHandler) MSG_SELL_ITEM_REQ(actorMsg *actor.ActorMessage) {

}

func (this *StorageHandler) MSG_GM_COMMAND_REQ(actorMsg *actor.ActorMessage) {
	gmReq := actorMsg.Proto().(*message.GmCommandReq)
	if flag.Lookup("gm").Value.(flag.Getter).Get().(int) > 0 {
		gmCommand(this.game.PlayerMgr().GetRoleBySession(actorMsg.GateSession()), gmReq.CMD)
	} else {
		log.Error("no open GM")
	}
	sendInfo := network.NewPbMessage(&message.GmCommandRes{Result: message.MSG_RESULT_SUCCESS}, message.MSG_GM_COMMAND_RES.Int32())
	system.Send(actorMsg.GateSession(), this.game.GetID(), actorMsg.SourceId(), sendInfo)
}

func (this *StorageHandler) UseItemRes(rid int64, actorMsg *actor.ActorMessage, result message.MSG_RESULT) {
	msg := &message.UseItemRes{
		Result: result,
	}
	LogResultHandler(rid, result)
	sendInfo := network.NewPbMessage(msg, message.MSG_USE_ITEM_RES.Int32())
	system.Send(actorMsg.GateSession(), this.game.GetID(), actorMsg.SourceId(), sendInfo)
}

func gmCommand(entity interfaces.IPlayer, command string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(string(debug.Stack()))
		}
	}()
	params := strings.Split(command, " ")
	log.Info("gm Command:" + command)
	f, ok := gmCommands[params[0]]
	if !ok {
		log.KVs(log.Fields{"command": params[0]}).Error("no this command")
		return
	}
	f(entity, params[1:])
}

func addItem(role interfaces.IPlayer, params []string) bool {
	itemId, _ := strconv.Atoi(params[0])
	itemNum, _ := strconv.Atoi(params[1])
	role.AddItemMap(map[int32]int64{int32(itemId): int64(itemNum)}, common.NewItemLog(common.CMD))
	return true
}
