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
var _CherryProduceMap = map[int64]*CherryProduce{}
var _CherryProduceArray = []*CherryProduce{}

type CherryProduce struct {
	data *_CherryProduce
}

// 类型结构
type _CherryProduce struct {
	INT_ID           int64 // INT_ID
	INT_RewardPoolId int64 // 产出道具池ID
	INT_Discount     int64 // 折算系数
}

func (c *CherryProduce) ID() int64           { return c.data.INT_ID }           //INT_ID
func (c *CherryProduce) RewardPoolId() int64 { return c.data.INT_RewardPoolId } //产出道具池ID
func (c *CherryProduce) Discount() int64     { return c.data.INT_Discount }     //折算系数

func HasCherryProduce(key int64) bool {
	_, ok := _CherryProduceMap[key]
	return ok
}

func GetCherryProduce(key int64) *CherryProduce {
	return _CherryProduceMap[key]
}

func RangeCherryProduce(fn func(i int, row *CherryProduce) (stop bool)) {
	for i, row := range _CherryProduceArray {
		if fn(i, row) {
			break
		}
	}
}

func LenCherryProduce() int { return len(_CherryProduceArray) }

func init() {
	loadfn["CherryProduce"] = loadCherryProduce
}

func loadCherryProduce(dir string) error {
	data, err := ioutil.ReadFile(path.Join(dir, "CherryProduce.json"))
	if err != nil {
		return fmt.Errorf("file=%v read err=%v", err.Error())
	}

	datas := []*_CherryProduce{}
	err = json.Unmarshal(data, &datas)
	if err != nil {
		return fmt.Errorf("file=%v parse err=%v", err.Error())
	}

	result_array := []*CherryProduce{}
	result_map := map[int64]*CherryProduce{}
	for _, row := range datas {
		data := &CherryProduce{data: row}
		result_array = append(result_array, data)
		result_map[row.INT_ID] = data
	}
	_CherryProduceArray = result_array
	_CherryProduceMap = result_map
	fmt.Printf("%-50v len:%v\n", "CherryProduce load finish! ", len(result_array))
	return nil
}
