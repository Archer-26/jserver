package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"reflect"
	"root/internal/common"

	//uuid "github.com/satori/go.uuid"
	"root/pkg/abtime"
	"root/pkg/log"
)

//玩家用户数据
type ModelUser struct {
	//ObjId            string       `gorm:"not null VARCHAR(32);primary_key;"`
	UID              int64        `gorm:"not null BIGINT(20);primary_key"` //角色ID
	Roles            RoleMap      `gorm:"type:json"`                       //RID=>RoleInfo 一个玩家可以拥有多个角色
	Passports        PassportsMap `gorm:"type:json"`                       //PassportType=>PID
	LastLogin        int64        `gorm:"BIGINT(20)"`                      //最近一次登录时间
	LastLoginRID     int64        `gorm:"BIGINT(20)"`                      //最近一次登录RID
	LastSecurityCode string       `gorm:"-"`                               //最近一次登录授权码
	ABTest           string       `gorm:"VARCHAR(64)"`                     //ab测试
	PackageCode      string       `gorm:"VARCHAR(64)"`                     //客户端包名
	CreateIP         string       `gorm:"VARCHAR(64)"`                     //绑定时的客户端IP
	CreatedAt        int64        `gorm:"BIGINT(20)"`                      //创建时间
	GateSession      string       `json:"-" gorm:"-"`                      //用户登录绑定的session
	GameActorId      string       `json:"-" gorm:"-"`                      //登录成功，分配的Game actor
}

func NewModelUser(uid int64, ip string) *ModelUser {
	result := &ModelUser{
		//ObjId:     uuid.NewV4().String(),
		UID:       uid,
		Roles:     RoleMap{},
		Passports: PassportsMap{},
		LastLogin: abtime.Now().Unix(),
		CreatedAt: abtime.Now().Unix(),
		CreateIP:  ip,
	}
	return result
}

// 加载所有user
func (this *ModelUser) LoadAll(db *gorm.DB) ([]*ModelUser, error) {
	all := []*ModelUser{}
	err := db.Find(&all).Error
	return all, err
}

func (this *ModelUser) ModelName() string {
	return reflect.TypeOf(this).Elem().Name()
}

//角色基础信息
type RoleInfo struct {
	RID       int64  //角色Id
	UID       int64  //游戏内账号ID
	Icon      string //头像
	Name      string //角色名
	SID       int64  //分配的区服ID
	LastLogin int64  //最近一次登录时间
	CreateAt  int64  //创建时间
	CreateIP  string //创建IP
	ABTest    string //ab测试
	IsDelete  bool   //删除标记
}

//第三方登录账号
type Passport struct {
	Type     string //账号类型
	PID      string //账号ID
	UID      int64  //游戏内账号ID
	CreateAt int64  //绑定时间
	CreateIP string //绑定时的客户端IP
}

func (this *ModelUser) String() string {
	jsonData, e := json.Marshal(this)
	if e != nil {
		log.KVs(log.Fields{"UID": this.UID, "error": e}).ErrorStack(3, "json.Marshal")
		return ""
	}
	return string(jsonData)
}

func (this *ModelUser) BelongDB() []string {
	return []string{common.LoginDbName}
}
