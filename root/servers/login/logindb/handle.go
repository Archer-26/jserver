package logindb

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
	"root/internal/common"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/expect"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/models"
)

const max_user_send_count = 5000

// login请求所有用户数据
func (this *LoginDB) INNER_MSG_L2D_ALL_USER_REQ(actorMsg *actor.ActorMessage) {
	// todo...暂时没有redis模块,数据都直接从mysql里取,先用闭包实现功能
	system.LocalSend(this.GetID(), this.mysql.GetID(), "INNER_MSG_L2D_ALL_USER_REQ", func() {
		allUser, err := (&models.ModelUser{}).LoadAll(this.mysql.DataBase[0])
		expect.Nil(err)

		spareCount := len(allUser)
		sendCount := 1
		for spareCount > max_user_send_count {
			bytes, err := msgpack.Marshal(allUser[:max_user_send_count])
			expect.Nil(err,log.Fields{"count(allUser)": len(allUser)})

			spareCount = spareCount - max_user_send_count
			allUser = allUser[max_user_send_count:]

			log.KV("size", len(bytes)).Info(colorized.Cyan("loginDB send partion of all user"))
			// 原则上，所有数据交互都通过db，不能直接mysqlActor传给login，所以这里再次丢给db
			sendMsg := network.NewPbMessage(&inner.D2LAllUserRes{Data: bytes}, inner.INNER_MSG_D2L_ALL_USER_RES.Int32())
			system.LocalSend(this.mysql.GetID(), this.GetID(), fmt.Sprintf("sendAllUser split:%v", sendCount), func() {
				system.Send("", this.GetID(), common.Login_Actor, sendMsg)
			})
			sendCount++
		}
		bytes, err := msgpack.Marshal(allUser)
		expect.Nil(err,log.Fields{"count(allUser)": len(allUser)})

		log.KV("size", len(bytes)).Info(colorized.Cyan("loginDB send partion of all user"))
		// 原则上，所有数据交互都通过db，不能直接从mysqlActor传给login，所以这里再次丢给db
		sendMsg := network.NewPbMessage(&inner.D2LAllUserRes{Data: bytes}, inner.INNER_MSG_D2L_ALL_USER_RES.Int32())
		system.LocalSend(this.mysql.GetID(), this.GetID(), fmt.Sprintf("sendAllUser split:%v", sendCount), func() {
			system.Send("", this.GetID(), common.Login_Actor, sendMsg)
			// 发个空消息，表示数据发送完成
			overSendMsg := network.NewPbMessage(&inner.D2LAllUserRes{Data: nil}, inner.INNER_MSG_D2L_ALL_USER_RES.Int32())
			system.Send("", this.GetID(), common.Login_Actor, overSendMsg)
		})
	})
}

// login 回存user
func (this *LoginDB) INNER_MSG_L2D_USER_SAVE(actorMsg *actor.ActorMessage) {
	// todo...暂时没有redis模块,数据都直接从mysql里取,先用闭包实现功能
	saveUser := actorMsg.Proto().(*inner.L2DUserSave)
	system.LocalSend(this.GetID(), this.mysql.GetID(), "INNER_MSG_ALL2D_MODEL_SAVE", func() {
		md := &models.ModelUser{}
		err := msgpack.Unmarshal(saveUser.Data, md)
		expect.Nil(err)

		this.mysql.DataBase[0].Save(md)
		log.KV("md", md.String()).Info("loginDB save new User")
	})
}
