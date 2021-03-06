package config_go

// Code generated by excelExoprt. DO NOT EDIT.
// source: CherrySlot.xlsx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// array and map
var _CherryWinTypeMap = map[int64]*CherryWinType{}
var _CherryWinTypeArray = []*CherryWinType{}

type CherryWinType struct {
	data *_CherryWinType
}

// 类型结构
type _CherryWinType struct {
	INT_ID              int64  // 赢钱称号
	INT_MIN             int64  // 倍率下限
	INT_MAX             int64  // 倍率上限
	STR_EFFECTS_ADDRESS string // 特效地址
}

func (c *CherryWinType) ID() int64               { return c.data.INT_ID }              //赢钱称号
func (c *CherryWinType) MIN() int64              { return c.data.INT_MIN }             //倍率下限
func (c *CherryWinType) MAX() int64              { return c.data.INT_MAX }             //倍率上限
func (c *CherryWinType) EFFECTS_ADDRESS() string { return c.data.STR_EFFECTS_ADDRESS } //特效地址

func HasCherryWinType(key int64) bool {
	_, ok := _CherryWinTypeMap[key]
	return ok
}

func GetCherryWinType(key int64) *CherryWinType {
	return _CherryWinTypeMap[key]
}

func RangeCherryWinType(fn func(i int, row *CherryWinType) (stop bool)) {
	for i, row := range _CherryWinTypeArray {
		if fn(i, row) {
			break
		}
	}
}

func LenCherryWinType() int { return len(_CherryWinTypeArray) }

func init() {
	loadfn["CherryWinType"] = loadCherryWinType
}

func loadCherryWinType(dir string) error {
	data, err := ioutil.ReadFile(path.Join(dir, "CherryWinType.json"))
	if err != nil {
		return fmt.Errorf("file=%v read err=%v", err.Error())
	}

	datas := []*_CherryWinType{}
	err = json.Unmarshal(data, &datas)
	if err != nil {
		return fmt.Errorf("file=%v parse err=%v", err.Error())
	}

	result_array := []*CherryWinType{}
	result_map := map[int64]*CherryWinType{}
	for _, row := range datas {
		data := &CherryWinType{data: row}
		result_array = append(result_array, data)
		result_map[row.INT_ID] = data
	}
	_CherryWinTypeArray = result_array
	_CherryWinTypeMap = result_map
	fmt.Printf("%-50v len:%v\n", "CherryWinType load finish! ", len(result_array))
	return nil
}
