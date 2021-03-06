package config_go

// Code generated by excelExoprt. DO NOT EDIT.
// source: BaseHero.xlsx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// array and map
var _BaseHeroMap = map[int64]*BaseHero{}
var _BaseHeroArray = []*BaseHero{}

type BaseHero struct {
	data *_BaseHero
}

// 类型结构
type _BaseHero struct {
	INT_ID      int64  // 英雄ID
	INT_QUALITY int64  // 品质_显示用
	STR_URL     string // 形象
}

func (c *BaseHero) ID() int64      { return c.data.INT_ID }      //英雄ID
func (c *BaseHero) QUALITY() int64 { return c.data.INT_QUALITY } //品质_显示用
func (c *BaseHero) URL() string    { return c.data.STR_URL }     //形象

func HasBaseHero(key int64) bool {
	_, ok := _BaseHeroMap[key]
	return ok
}

func GetBaseHero(key int64) *BaseHero {
	return _BaseHeroMap[key]
}

func RangeBaseHero(fn func(i int, row *BaseHero) (stop bool)) {
	for i, row := range _BaseHeroArray {
		if fn(i, row) {
			break
		}
	}
}

func LenBaseHero() int { return len(_BaseHeroArray) }

func init() {
	loadfn["BaseHero"] = loadBaseHero
}

func loadBaseHero(dir string) error {
	data, err := ioutil.ReadFile(path.Join(dir, "BaseHero.json"))
	if err != nil {
		return fmt.Errorf("file=%v read err=%v", err.Error())
	}

	datas := []*_BaseHero{}
	err = json.Unmarshal(data, &datas)
	if err != nil {
		return fmt.Errorf("file=%v parse err=%v", err.Error())
	}

	result_array := []*BaseHero{}
	result_map := map[int64]*BaseHero{}
	for _, row := range datas {
		data := &BaseHero{data: row}
		result_array = append(result_array, data)
		result_map[row.INT_ID] = data
	}
	_BaseHeroArray = result_array
	_BaseHeroMap = result_map
	fmt.Printf("%-50v len:%v\n", "BaseHero load finish! ", len(result_array))
	return nil
}
