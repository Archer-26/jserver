package models

import (
	"database/sql/driver"
	"encoding/json"
)

/*
	mysql5.7 版本支持json格式
	为json定义新的类型并实现 Value() (driver.Value, error) 和 Scan(input interface{}) error
	注:
		1.Value方法接收者必须是值类型！！	 (否则增删查改会报错)
		2.Scan方法接收者必须是指针类型！！   (否则增删查改会报错)
    example:
		gorm的 mssql.JSON 定义
*/
type (
	RoleMap      map[int64]*RoleInfo
	PassportsMap map[string]*Passport
	ItemMap      map[int32]*Item
	BuildingMap  map[int32]*Building
	IslandMap    map[int32]*Island
)

// RoleMap
func (this RoleMap) Value() (driver.Value, error) {
	return json.Marshal(this)
}

func (this *RoleMap) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), this)
}

// PassportsMap
func (this PassportsMap) Value() (driver.Value, error) {
	return json.Marshal(this)
}

func (this *PassportsMap) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), this)
}

// ItemMap
func (this ItemMap) Value() (driver.Value, error) {
	return json.Marshal(this)
}

func (this *ItemMap) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), this)
}

// BuildingMap
func (this BuildingMap) Value() (driver.Value, error) {
	return json.Marshal(this)
}

func (this *BuildingMap) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), this)
}

// Island
func (this IslandMap) Value() (driver.Value, error) {
	return json.Marshal(this)
}

func (this *IslandMap) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), this)
}
