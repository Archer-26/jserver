package island

import (
	"root/internal/common"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
	"root/servers/internal/models"
	"root/servers/internal/mysql/i_mysql"
)

type (
	Island struct {
		RID  int64
		game interfaces.IGame
		*models.ModelIsland
	}
)

func (this *Island) Init(rid int64, g interfaces.IGame) {
	this.RID = rid
	this.game = g
	this.ModelIsland = models.NewModelIsland(rid)
}

func (this *Island) Sync() {
	islands := []*message.Island{}
	for _, v := range this.IslandList {
		islands = append(islands, &message.Island{
			Id:           v.IslandId,
			CheckpointId: v.CheckpointId,
			BuildingMap:  v.PlaceMap,
			JsonInfo: v.JsonInfo,
		})
	}
	islandInfo := &message.NotifyIslandInfo{
		Islands: islands,
	}

	msg := network.NewPbMessage(islandInfo, message.MSG_NOTIFY_ISLAND_INFO.Int32())
	this.player().SendMsg(msg)
}

func (this *Island) Save() {
	this.game.Cache().Save(this.RID, this.ModelIsland)
}

func (this *Island) Place(island *message.Island) {
	if _, ok := this.ModelIsland.IslandList[island.Id]; !ok {
		this.ModelIsland.IslandList[island.Id] = &models.Island{
			IslandId:     island.Id,
			CheckpointId: island.CheckpointId,
			PlaceMap:     island.BuildingMap,
		}
	} else {
		this.ModelIsland.IslandList[island.Id].CheckpointId = island.CheckpointId
		this.ModelIsland.IslandList[island.Id].PlaceMap = island.BuildingMap
		this.ModelIsland.IslandList[island.Id].JsonInfo = island.JsonInfo
	}

	log.KVs(log.Fields{"rid": this.RID,"island":island.String()}).Debug("Island Place")
	this.game.Cache().Save(this.RID, this.ModelIsland)
}

func (this *Island) CheckpointFinish(island *message.CheckpointFinishReq) {
	this.player().AddGold(island.Gold, common.NewItemLog(common.CHECKPOINT_FINISH))
	if island.Success == 0 {
		if _, ok := this.ModelIsland.IslandList[int32(island.IslandId)]; !ok {
			this.ModelIsland.IslandList[int32(island.IslandId)] = &models.Island{
				IslandId:     int32(island.IslandId),
				CheckpointId: int32(island.CheckpointId),
				PlaceMap:     map[int32]int32{},
			}
		} else {
			this.ModelIsland.IslandList[int32(island.IslandId)].CheckpointId = int32(island.CheckpointId)
		}

	}
	log.KVs(log.Fields{"rid": this.RID, "success": island.Success, "gold": island.Gold, "pointId": island.CheckpointId}).
		Debug("Island CheckpointFinish")
	this.game.Cache().Save(this.RID, this.ModelIsland)
}
func (this *Island) SaveJsonInfo(island *message.SaveIslandInfoReq) {
	if _, ok := this.ModelIsland.IslandList[int32(island.IslandId)]; !ok {
		this.ModelIsland.IslandList[int32(island.IslandId)] = &models.Island{
			IslandId: int32(island.IslandId),
			JsonInfo: island.JsonInfo,
		}
	} else {
		this.ModelIsland.IslandList[int32(island.IslandId)].JsonInfo = island.JsonInfo
	}
	log.KVs(log.Fields{"rid": this.RID, "island": island.String()}).Debug("Island SaveJsonInfo")
	this.game.Cache().Save(this.RID, this.ModelIsland)
}

func (this *Island) Model() i_mysql.IDbTable {
	return this.ModelIsland
}

func (this *Island) player() interfaces.IPlayer {
	return this.game.PlayerMgr().GetRoleByRID(this.RID)
}
