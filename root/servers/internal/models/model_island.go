package models

import (
	"encoding/json"
	"reflect"
	"root/internal/common"
	"root/pkg/log"
)

// game内的玩家对象
type ModelIsland struct {
	RID        int64     `gorm:"not null; default 0 BIGINT(20); primary_key"` //角色ID
	IslandList IslandMap `gorm:"type:json"`                                   //完成的所有岛屿
}

// 道具基础信息
type Island struct {
	IslandId     int32           `json:"IslandId"`     // 岛屿Id
	CheckpointId int32           `json:"CheckpointId"` // 关卡Id
	PlaceMap     map[int32]int32 `json:"PlaceMap"`
	JsonInfo     string          `json:"JsonInfo"`
}

func NewModelIsland(rid int64) *ModelIsland {
	return &ModelIsland{
		RID:        rid,
		IslandList: IslandMap{},
	}
}

func (this *ModelIsland) ModelName() string {
	return reflect.TypeOf(this).Elem().Name()
}

func (this *ModelIsland) String() string {
	jsonData, e := json.Marshal(this)
	if e != nil {
		log.KVs(log.Fields{"RID": this.RID, "error": e}).ErrorStack(3, "json.Marshal")
		return ""
	}
	return string(jsonData)
}

func (this *ModelIsland) BelongDB() []string {
	return []string{common.GameDbName}
}
