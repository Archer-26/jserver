package models

import (
	"encoding/json"
	"reflect"
	"root/internal/common"
	"root/pkg/log"
)

// 老虎机相关的角色数据
type ModelSlot struct {
	RID      int64 `gorm:"not null BIGINT(20);primary_key"` // 角色ID
	IsSelect bool  `gorm:"not null;"`                       // 是否触发Scatter
	LastBet  int64 `gorm:"not null; BIGINT(20);"`           // 最后一次押注等级
}

func NewModelSlot(rid int64) *ModelSlot {
	return &ModelSlot{
		RID:      rid,
		IsSelect: false,
		LastBet:  0,
	}
}

func (this *ModelSlot) ModelName() string {
	return reflect.TypeOf(this).Elem().Name()
}

func (this *ModelSlot) String() string {
	jsonData, e := json.Marshal(this)
	if e != nil {
		log.KVs(log.Fields{"RID": this.RID, "error": e}).ErrorStack(3, "json.Marshal")
		return ""
	}
	return string(jsonData)
}

func (this *ModelSlot) BelongDB() []string {
	return []string{common.GameDbName}
}
