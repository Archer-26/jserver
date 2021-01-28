package models

import (
	"encoding/json"
	"reflect"
	"root/internal/common"
	"root/pkg/log"
)

// 骰子游戏信息
type ModelDice struct {
	RID       int64     `gorm:"not null; default 0 BIGINT(20); primary_key"` //角色ID
	GameInfo string 	`gorm:"type:MEDIUMTEXT"`//保存的游戏信息
}


func NewModelDice(rid int64) *ModelDice {
	return &ModelDice{
		RID:       rid,
		GameInfo: "",
	}
}

func (this *ModelDice) ModelName() string {
	return reflect.TypeOf(this).Elem().Name()
}

func (this *ModelDice) String() string {
	jsonData, e := json.Marshal(this)
	if e != nil {
		log.KVs(log.Fields{"RID": this.RID, "error": e}).ErrorStack(3, "json.Marshal")
		return ""
	}
	return string(jsonData)
}

func (this *ModelDice) BelongDB() []string {
	return []string{common.GameDbName}
}
