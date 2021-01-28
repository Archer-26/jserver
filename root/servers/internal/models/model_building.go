package models

import (
	"encoding/json"
	"reflect"
	"root/internal/common"
	"root/internal/config/config_global"
	"root/pkg/log"
	"root/servers/internal/message"
)

type ModelBuildings struct {
	RID       int64       `gorm:"not null BIGINT(20);primary_key"` // 角色ID
	Buildings BuildingMap `gorm:"type:json"`                       // 所有建筑
}

func NewModelBuildings(rid int64) *ModelBuildings {
	return &ModelBuildings{
		RID:       rid,
		Buildings: BuildingMap{},
	}
}

func (this *ModelBuildings) ModelName() string {
	return reflect.TypeOf(this).Elem().Name()
}

func (this *ModelBuildings) String() string {
	jsonData, e := json.Marshal(this)
	if e != nil {
		log.KVs(log.Fields{"RID": this.RID, "error": e}).ErrorStack(3, "json.Marshal")
		return ""
	}
	return string(jsonData)
}

func (this *ModelBuildings) BelongDB() []string {
	return []string{common.GameDbName}
}

// 建筑基础信息
type Building struct {
	Id         int32           `json:"Id"`
	Lv         int32           `json:"Lv"`
	Star       int32           `json:"Star"`
	PropsLevel map[int32]int32 `json:"PropsLevel"`
}

func NewBuilding(buildingId, lv, star int32) *Building {
	building := &Building{
		Id:         buildingId,
		Lv:         lv,
		Star:       star,
		PropsLevel: make(map[int32]int32),
	}
	for attr, _ := range config_global.BuildingWithProp[int64(buildingId)] {
		building.PropsLevel[int32(attr)] = 1
	}
	return building
}

func (building *Building) SetLevel(newLevel int32) {
	building.Lv = newLevel
}

func (building *Building) SetStar(newStar int32) {
	building.Star = newStar
}

func (building *Building) SetPropLevel(propId int32, newLevel int32) {
	building.PropsLevel[propId] = newLevel
}

func (building *Building) ToMessage() *message.Building {
	return &message.Building{
		Id:         building.Id,
		Level:      building.Lv,
		Star:       building.Star,
		PropsLevel: building.PropsLevel,
	}
}
