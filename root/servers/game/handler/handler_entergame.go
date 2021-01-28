package handler

import (
	"root/internal/common"
	"root/internal/system"
	"root/pkg/abtime"
	"root/pkg/actor"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
)

type enterGameHandler struct {
	game interfaces.IGame
}

func InitEnterGameHandler(g interfaces.IGame) {
	handler := &enterGameHandler{
		game: g,
	}
	g.RegistMsg(message.MSG_ENTER_GAME_REQ.Int32(), handler.MSG_ENTER_GAME_REQ)
}

// 玩家请求进入游戏
func (this *enterGameHandler) MSG_ENTER_GAME_REQ(actorMsg *actor.ActorMessage) {
	enterGameReq := actorMsg.Proto().(*message.EnterGameReq)
	logField := log.Fields{"NewPlayer": enterGameReq.NewPlayer, "RID": enterGameReq.RID, "UID": enterGameReq.UID, "actorMsg.gateSession()": actorMsg.GateSession(), "source": actorMsg.SourceId()}
	if !common.IsLogin(enterGameReq.RID, enterGameReq.UID, enterGameReq.SID, enterGameReq.LoginTime, enterGameReq.NewPlayer, enterGameReq.SecurityCode) {
		log.KVs(logField).Warn("login auth failed")
		this.RespEnterGame(actorMsg.GateSession(), message.MSG_RESULT_FAILED, actorMsg.SourceId(), nil)
		return
	}

	if enterGameReq.NewPlayer {
		role := this.game.PlayerMgr().GetRoleByRID(enterGameReq.RID)
		if role != nil {
			log.KV("RID", role.RID()).Warn("player was not found")
		} else {
			role = this.game.PlayerMgr().NewRole(actorMsg.GateSession(), enterGameReq.RID, enterGameReq.UID, enterGameReq.SID, enterGameReq.Passport, enterGameReq.PassportType) // 新玩家
		}
		this.RespEnterGame(actorMsg.GateSession(), message.MSG_RESULT_SUCCESS, actorMsg.SourceId(), role)
	} else {
		role := this.game.PlayerMgr().GetRoleByRID(enterGameReq.RID)
		if role != nil {
			this.game.PlayerMgr().RelateSession2Role(actorMsg.GateSession(), role)
			role.BindSession(actorMsg.GateSession())
			this.RespEnterGame(actorMsg.GateSession(), message.MSG_RESULT_SUCCESS, actorMsg.SourceId(), role)
		} else {
			// 防止sync.pool复用actorMsg，导致闭包数据错误,需要的数据提出来用
			gateSession := actorMsg.GateSession()
			sourceId := actorMsg.SourceId()
			this.game.PlayerMgr().GetRoleByAsync(enterGameReq.RID, func(mod interfaces.WholeRoleBytes) { // 本地没有角色数据，就找db要
				var callback_role interfaces.IPlayer
				callback_role = this.game.PlayerMgr().GetRoleByRID(enterGameReq.RID)
				if callback_role == nil {
					if mod == nil {
						// 触发这里的情况：1.要么客户端NewPlayer字段处理错误，2.要么第一次登录成功，但是没有正常进入游戏，导致login有角色，但game没有
						log.KVs(logField).KV("rid", enterGameReq.RID).Warn("mysql was not found this player")
						callback_role = this.game.PlayerMgr().NewRole(gateSession, enterGameReq.RID, enterGameReq.UID, enterGameReq.SID, enterGameReq.Passport, enterGameReq.PassportType) // 异常触发新玩家
					} else {
						callback_role = this.game.PlayerMgr().RoleMod(gateSession, enterGameReq.RID, mod)
					}
				} else {
					log.KV("RID", enterGameReq.RID).Warn(colorized.Red("GetRoleByAsync so slowly!!! client repeated login!"))
				}

				log.KVs(logField).Info("GetRoleByAsync callback")
				this.RespEnterGame(gateSession, message.MSG_RESULT_SUCCESS, sourceId, callback_role)
			})
		}
	}
	log.KVs(logField).Info("player enter game")
}

func (this *enterGameHandler) RespEnterGame(gateSession string, result message.MSG_RESULT, sourceId string, role interfaces.IPlayer) {
	msg := &message.EnterGameRes{
		Result:     result,
		ServerTime: abtime.Now().UnixNano(),
		ServerZone: int32(abtime.TimeZone / 3600),
	}

	sendInfo := network.NewPbMessage(msg, message.MSG_ENTER_GAME_RES.Int32())
	system.Send(gateSession, this.game.GetID(), sourceId, sendInfo)
	if result == message.MSG_RESULT_SUCCESS {
		role.SyncAllData()
	}
}
