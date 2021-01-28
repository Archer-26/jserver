package models

import (
	"encoding/json"
	"reflect"
	"root/internal/common"
	"root/pkg/log"
)

type ModelItems struct {
	RID      int64   `gorm:"not null BIGINT(20);primary_key"` // 角色ID
	ItemList ItemMap `gorm:"type:json"`                       // 所有道具
}

// 道具基础信息
type Item struct {
	ItemId     int32 `json:"ItemId"`
	ItemCount  int64 `json:"ItemCount"`
	CreateTime int64 `json:"CreateTime"` //首次创建时间 ;用来做过期
}

func NewModelItems(rid int64) *ModelItems {
	return &ModelItems{
		RID:      rid,
		ItemList: ItemMap{},
	}
}

func (this *ModelItems) ModelName() string {
	return reflect.TypeOf(this).Elem().Name()
}

func (this *ModelItems) String() string {
	jsonData, e := json.Marshal(this)
	if e != nil {
		log.KVs(log.Fields{"RID": this.RID, "error": e}).ErrorStack(3, "json.Marshal")
		return ""
	}
	return string(jsonData)
}

func (this *ModelItems) BelongDB() []string {
	return []string{common.GameDbName}
}
