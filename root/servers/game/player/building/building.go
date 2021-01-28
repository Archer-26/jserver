package building

import (
	"root/internal/common"
	"root/pkg/ev"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/game/event"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
	"root/servers/internal/models"
	"root/servers/internal/mysql/i_mysql"
)

type (
	Building struct {
		RID  int64
		game interfaces.IGame
		*models.ModelBuildings
	}
)

func (this *Building) Init(rid int64, g interfaces.IGame) {
	this.RID = rid
	this.game = g
	this.ModelBuildings = models.NewModelBuildings(rid)
}

func (this *Building) Sync() {
	var buildings []*message.Building
	for _, v := range this.Buildings {
		buildings = append(buildings, v.ToMessage())
	}
	this.notifyChange(buildings)
}

func (this *Building) AddBuilding(buildingId int32, level int32, star int32, itemLog *common.ItemLog) {
	field := log.Fields{"rid": this.RID, "buildingId": buildingId, "level": level, "star": star, "itemLog": itemLog}
	_, ok := this.Buildings[buildingId]
	if ok {
		//todo 变成道具？
		log.KVs(field).Info("modify to other")
		return
	}
	building := models.NewBuilding(buildingId, level, star)
	this.Buildings[buildingId] = building
	this.Save()
	this.notifyChange([]*message.Building{building.ToMessage()})
	log.KVs(field).Info("AddBuilding")
}

func (this *Building) SetBuildingLevel(buildingId int32, newLevel int32, itemLog *common.ItemLog) {
	field := log.Fields{"rid": this.RID, "buildingId": buildingId, "newLevel": newLevel, "itemLog": itemLog}
	building, ok := this.Buildings[buildingId]
	if !ok {
		log.KVs(field).Error("no this buildingId")
		return
	}
	field.AddFiled("oldLevel", building.Lv)
	building.SetLevel(newLevel)
	this.Save()
	this.notifyChange([]*message.Building{building.ToMessage()})
	log.KVs(field).Info("SetBuildingLevel")
}

func (this *Building) SetBuildingStar(buildingId int32, newStar int32, itemLog *common.ItemLog) {
	field := log.Fields{"rid": this.RID, "buildingId": buildingId, "newStar": newStar, "itemLog": itemLog}
	building, ok := this.Buildings[buildingId]
	if !ok {
		log.KVs(field).Error("no this buildingId")
		return
	}
	field.AddFiled("oldStar", building.Star)
	building.SetStar(newStar)
	this.Save()
	this.notifyChange([]*message.Building{building.ToMessage()})
	log.KVs(field).Info("SetBuildingStar")
}

func (this *Building) SetBuildingPropLevel(buildingId int32, propId int32, newPropLevel int32, itemLog *common.ItemLog) {
	field := log.Fields{"rid": this.RID, "buildingId": buildingId, "propId": propId, "newPropLevel": newPropLevel, "itemLog": itemLog}
	building, ok := this.Buildings[buildingId]
	if !ok {
		log.KVs(field).Error("no this buildingId")
		return
	}
	field.AddFiled("oldPropLevel", building.PropsLevel[propId])
	building.SetPropLevel(propId, newPropLevel)
	this.Save()
	this.notifyChange([]*message.Building{building.ToMessage()})
	log.KVs(field).Info("SetBuildingStar")
}

func (this *Building) OnEvent(evt ev.IEvent) {
	switch evt.EType() {
	case event.EV_ROLE_UPGRADE:
		//evet := evt.(*event.Ev_Role_Upgrade)
		//// Level todo ..bla ......

	}
}

func (this *Building) GetBuilding(buildingId int32) *models.Building {
	return this.Buildings[buildingId]
}
func (this *Building) Save() {
	this.game.Cache().Save(this.RID, this.ModelBuildings)
}

func (this *Building) Model() i_mysql.IDbTable {
	return this.ModelBuildings
}

func (this *Building) Player() interfaces.IPlayer {
	return this.game.PlayerMgr().GetRoleByRID(this.RID)
}

func (this *Building) notifyChange(buildings []*message.Building) {
	buildingInfo := &message.NotifyBuildingInfo{Buildings: buildings}
	msg := network.NewPbMessage(buildingInfo, message.MSG_NOTIFY_BUILDING_INFO.Int32())
	this.Player().SendMsg(msg)
}
