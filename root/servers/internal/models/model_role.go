package models

import (
	"encoding/json"
	"reflect"
	"root/internal/common"
	"root/pkg/log"
)

// game内的玩家对象
type ModelRole struct {
	//ObjId        string `gorm:"not null; VARCHAR(32);primary_key"`
	RID          int64  `gorm:"not null; default 0 BIGINT(20); primary_key"` //角色ID
	UID          int64  `gorm:"not null; BIGINT(20);index"`                  //userID
	SID          int64  `gorm:"not null; BIGINT(20);"`                       //sessionID
	Passport     string `gorm:"VARCHAR(128)"`
	PassportType string `gorm:"VARCHAR(128)"`
	Name         string `gorm:"VARCHAR(64); index"`
	Level        int32  `gorm:"not null; default 0 INT(11)"`
	Exp          int64  `gorm:"not null; default 0 BIGINT(20)"`
	Gold         int64  `gorm:"VARCHAR(64); index"`
	Country      string `gorm:"not null; default '' VARCHAR(32)"`
	CreatedAt    int64  `gorm:"BIGINT(20)"`
	IsDelete     bool   `gorm:"TINYINT(1)"` //是否删除
}

func NewModelRole(rid int64) *ModelRole {
	return &ModelRole{
		RID:   rid,
		Level: 1,
		Exp:   0,
	}
}

func (this *ModelRole) ModelName() string {
	return reflect.TypeOf(this).Elem().Name()
}

// * 所有表数据存储统一由 INNER_MSG_ALL2D_MODEL_SAVE 方法处理
//func (this *ModelRole) Save(mysql *gorm.DB) error {
//	return mysql.Save(this).Error
//}

func (this *ModelRole) String() string {
	jsonData, e := json.Marshal(this)
	if e != nil {
		log.KVs(log.Fields{"RID": this.RID, "error": e}).ErrorStack(3, "json.Marshal")
		return ""
	}
	return string(jsonData)
}

func (this *ModelRole) BelongDB() []string {
	return []string{common.GameDbName}
}
