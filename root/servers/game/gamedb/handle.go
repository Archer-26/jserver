package gamedb

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/pkg/tools/structenh"
	"root/servers/game/interfaces"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/models"
	"root/servers/internal/mysql/i_mysql"
)

// game请求角色数据
func (this *GameDB) INNER_MSG_G2D_ROLE_REQ(actorMsg *actor.ActorMessage) {
	/*
		防止actorMsg被pool回收后，对象复用，导致闭包内数据错乱，所以先把actorMsg数据取出来
	*/
	roleReq := actorMsg.Proto().(*inner.G2DRoleReq)
	session := actorMsg.GateSession()
	sourceId := actorMsg.SourceId()
	// todo...暂时没有redis模块,数据都直接从mysql里取,先用闭包实现功能
	system.LocalSend(this.GetID(), this.mysql.GetID(), "INNER_MSG_G2D_ROLE_REQ", func() {
		totalData, valid := this.g2DGetRole(roleReq.RID, _hashIndex(roleReq.RID))
		sendMsg := network.NewPbMessage(&inner.D2GRoleRes{RID: roleReq.RID, CallbackId: roleReq.CallbackId, WholeInfo: totalData, Valid: valid}, inner.INNER_MSG_D2G_ROLE_RES.Int32())
		system.LocalSend(this.mysql.GetID(), this.GetID(), fmt.Sprintf("send player callbackid:%v rid:%v", roleReq.CallbackId, roleReq.RID), func() {
			system.Send(session, this.GetID(), sourceId, sendMsg)
		})
	})
}

// 存储数据
func (this *GameDB) INNER_MSG_ALL2D_MODEL_SAVE(actorMsg *actor.ActorMessage) {
	model := actorMsg.Proto().(*inner.G2DModelSave)
	// todo...存数据暂时直接丢给mysql存储
	system.LocalSend(this.GetID(), this.mysql.GetID(), "INNER_MSG_ALL2D_MODEL_SAVE", func() {
		this.all2DModelSave(_hashIndex(model.RID), model.ModelName, model.Data)
	})
}

// game会在cachemod发送完成后，发送此消息，表示回存数据发送完成
func (this *GameDB) INNER_MSG_G2D_GAME_STOP(actorMsg *actor.ActorMessage) {
	//model := actorMsg.Proto().(*inner_message.G2DGameStop)
	system.LocalSend(this.GetID(), this.mysql.GetID(), "INNER_MSG_G2D_GAME_STOP", func() {
		log.Info(colorized.Red("finished stored data"))
		this.mysql.SuspendStop()
		system.LocalSend(this.mysql.GetID(), this.GetID(), "INNER_MSG_G2D_GAME_STOP", func() {
			this.SuspendStop()
		})
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func _hashIndex(rid int64) int {
	index := rid / 100
	return int(index % databaseCount)
}

func (this *GameDB) g2DGetRole(rid int64, dbIndex int) (interfaces.WholeRoleBytes, int32) {
	totalData := make(interfaces.WholeRoleBytes)
	var bytes []byte
	valid := int32(0)

	modRole := &models.ModelRole{}
	tabName := (modRole).ModelName()
	// 先找role表，role表里没有无需再查其他表
	clone := structenh.DeepClone(modRole).(i_mysql.IDbTable)
	err := this.mysql.DataBase[dbIndex].Find(clone, rid).Error
	if err != nil {
		log.KVs(log.Fields{"mod": tabName, "err": err}).Yellow().Warn("Load player info faild")
		valid = 1
		return totalData, valid
	} else {
		bytes, valid = _msgbytes(clone, rid)
		totalData[tabName] = bytes
	}

	// 玩家其他模块数据
	i_mysql.RangeTableWithDbName(this.mysql.DataBaseName, func(iDbTable i_mysql.IDbTable) bool {
		if tabName != iDbTable.ModelName() {
			clone = structenh.DeepClone(iDbTable).(i_mysql.IDbTable)
			err = this.mysql.DataBase[dbIndex].Find(clone, rid).Error
			if err == nil {
				bytes, valid = _msgbytes(clone, rid)
				totalData[iDbTable.ModelName()] = bytes
			} else {
				log.KVs(log.Fields{"mod": iDbTable.ModelName(), "err": err}).Yellow().Warn("mysql find error")
			}
		}
		return false
	})
	return totalData, valid
}

func (this *GameDB) all2DModelSave(dataBaseIndex int, modelName string, data []byte) {
	lfd := log.Fields{"modName": modelName, "dataBaseIndex": dataBaseIndex}
	dbTable := i_mysql.GetModelByName(modelName)
	if dbTable == nil {
		log.KVs(lfd).Error("not found model in typeOfMap")
		return
	}
	clone := structenh.DeepClone(dbTable).(i_mysql.IDbTable)
	err := msgpack.Unmarshal(data, clone)
	if err != nil {
		log.KVs(lfd).KV("error", err).Error("INNER_MSG_ALL2D_MODEL_SAVE msgpack.Unmarshal")
		return
	}
	this.mysql.DataBase[dataBaseIndex].Save(clone)
	log.KVs(lfd).KV("data", clone).Green().Info("GameDb save mod data")
}

func _msgbytes(v interface{}, RID int64) ([]byte, int32) {
	bytes, err := msgpack.Marshal(v)
	if err != nil {
		log.KVs(log.Fields{"err": err, "RID": RID}).ErrorStack(3, "msgpack.Marshal error")
		return nil, 1
	}
	return bytes, 0
}
