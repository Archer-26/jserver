package login

import (
	"root/internal/common"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/message"
)

// 登录
func (this *Login) MSG_TYPE_LOGIN_REQ(actorMsg *actor.ActorMessage) {
	loginReq := actorMsg.Proto().(*message.LoginReq)
	role, user, isNewUser := UserMgr.Login(loginReq.Passport, loginReq.PassportType, loginReq.ClientIP, loginReq.PackageCode)

	if isNewUser {
		// todo ??????????
	}
	// 取到user分配gameActor
	if user.GameActorId == "" {
		gameActor := GameMgr.GetAvailableGame()
		if gameActor == "" {
			log.Warn(colorized.Red("without available game actor"))
			loginResp := &message.LoginRes{Result: message.MSG_RESULT_SECURITYCODE_CHECK_FAILED}
			msg_res := network.NewPbMessage(loginResp, message.MSG_LOGIN_RES.Int32())
			system.Send(actorMsg.GateSession(), this.GetID(), actorMsg.SourceId(), msg_res)
			return
		} else {
			user.GameActorId = gameActor
		}
	}

	// 新旧链接相同，不需要处理旧链接
	if user.GateSession != "" && user.GateSession != actorMsg.GateSession() {
		msgDisabled := network.NewPbMessage(&inner.L2GTUserSessionDisabled{GateSession: user.GateSession, UID: user.UID}, inner.INNER_MSG_L2GT_USER_SESSION_DISABLED.Int32())
		atrId, _ := common.SplitGateSession(user.GateSession)
		system.Send(actorMsg.GateSession(), this.GetID(), atrId, msgDisabled) // 通知旧的gate，断开用户旧链接
		log.KVs(log.Fields{"UID": user.UID, "gateSession": user.GateSession}).Warn("user request login present of old gateSession")
	}
	user.GateSession = actorMsg.GateSession()
	// 先通知gate绑定gameActor
	l2gMsg := network.NewPbMessage(&inner.L2GTSessionAssignGame{
		GateSession: actorMsg.GateSession(),
		GameActorId: user.GameActorId,
	}, inner.INNER_MSG_L2GT_SESSION_ASSIGN_GAME.Int32())
	system.Send(actorMsg.GateSession(), this.GetID(), actorMsg.SourceId(), l2gMsg)

	// 登录成功，通知用户，下一步客户端请求进入游戏
	loginResp := network.NewPbMessage(&message.LoginRes{
		Result:         message.MSG_RESULT_SUCCESS,
		RID:            role.RID,
		UID:            user.UID,
		SID:            role.SID,
		Timestamp:      role.LastLogin,
		SecurityCode:   user.LastSecurityCode,
		UserCreateTime: user.CreatedAt,
		RoleCreateTime: role.CreateAt,
		NewPlayer:      isNewUser,
	}, message.MSG_LOGIN_RES.Int32())
	system.Send(actorMsg.GateSession(), this.GetID(), actorMsg.SourceId(), loginResp)
	log.KVs(log.Fields{"UID": user.UID, "RID": role.RID, "SID": role.SID, "gameActor": user.GameActorId, "isNew": isNewUser, "gateSession": user.GateSession}).
		Info(colorized.Yellow("login success"))
}
