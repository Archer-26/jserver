package config_go

// Code generated by excelExoprt. DO NOT EDIT.
// source: BasePveStage.xlsx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// array and map
var _BasePveStageMap = map[int64]*BasePveStage{}
var _BasePveStageArray = []*BasePveStage{}

type BasePveStage struct {
	data *_BasePveStage
}

// 类型结构
type _BasePveStage struct {
	INT_ID                   int64     // 关卡ID
	STR_NAME                 string    // 本地化
	INT_NEXT                 int64     // 下一关
	INT_ISLAND_ID            int64     // 归属岛屿
	ARRAYINT_WAVE_ID         array_int // 怪物
	ARRAYINT_WAVE_START_TIME array_int // 每波出发时间(ms)-波与波之间的间隔时间
	INT_VICTORY_POINT        int64     // 通关需求积分
	INT_REWARD_POOL          int64     // 奖励
	INT_IDLE_REWARD          int64     // 离线奖励/min
}

func (c *BasePveStage) ID() int64        { return c.data.INT_ID }        //关卡ID
func (c *BasePveStage) NAME() string     { return c.data.STR_NAME }      //本地化
func (c *BasePveStage) NEXT() int64      { return c.data.INT_NEXT }      //下一关
func (c *BasePveStage) ISLAND_ID() int64 { return c.data.INT_ISLAND_ID } //归属岛屿
//怪物
func (c *BasePveStage) Len_WAVE_ID() int          { return c.data.ARRAYINT_WAVE_ID.Len() }
func (c *BasePveStage) Get_WAVE_ID(key int) int64 { return c.data.ARRAYINT_WAVE_ID.Get(key) }
func (c *BasePveStage) Range_WAVE_ID(fn func(int, int64) (stop bool)) {
	c.data.ARRAYINT_WAVE_ID.Range(fn)
}
func (c *BasePveStage) Copy_WAVE_ID() array_int { return c.data.ARRAYINT_WAVE_ID.Copy() }

//每波出发时间(ms)-波与波之间的间隔时间
func (c *BasePveStage) Len_WAVE_START_TIME() int { return c.data.ARRAYINT_WAVE_START_TIME.Len() }
func (c *BasePveStage) Get_WAVE_START_TIME(key int) int64 {
	return c.data.ARRAYINT_WAVE_START_TIME.Get(key)
}
func (c *BasePveStage) Range_WAVE_START_TIME(fn func(int, int64) (stop bool)) {
	c.data.ARRAYINT_WAVE_START_TIME.Range(fn)
}
func (c *BasePveStage) Copy_WAVE_START_TIME() array_int {
	return c.data.ARRAYINT_WAVE_START_TIME.Copy()
}
func (c *BasePveStage) VICTORY_POINT() int64 { return c.data.INT_VICTORY_POINT } //通关需求积分
func (c *BasePveStage) REWARD_POOL() int64   { return c.data.INT_REWARD_POOL }   //奖励
func (c *BasePveStage) IDLE_REWARD() int64   { return c.data.INT_IDLE_REWARD }   //离线奖励/min

func HasBasePveStage(key int64) bool {
	_, ok := _BasePveStageMap[key]
	return ok
}

func GetBasePveStage(key int64) *BasePveStage {
	return _BasePveStageMap[key]
}

func RangeBasePveStage(fn func(i int, row *BasePveStage) (stop bool)) {
	for i, row := range _BasePveStageArray {
		if fn(i, row) {
			break
		}
	}
}

func LenBasePveStage() int { return len(_BasePveStageArray) }

func init() {
	loadfn["BasePveStage"] = loadBasePveStage
}

func loadBasePveStage(dir string) error {
	data, err := ioutil.ReadFile(path.Join(dir, "BasePveStage.json"))
	if err != nil {
		return fmt.Errorf("file=%v read err=%v", err.Error())
	}

	datas := []*_BasePveStage{}
	err = json.Unmarshal(data, &datas)
	if err != nil {
		return fmt.Errorf("file=%v parse err=%v", err.Error())
	}

	result_array := []*BasePveStage{}
	result_map := map[int64]*BasePveStage{}
	for _, row := range datas {
		data := &BasePveStage{data: row}
		result_array = append(result_array, data)
		result_map[row.INT_ID] = data
	}
	_BasePveStageArray = result_array
	_BasePveStageMap = result_map
	fmt.Printf("%-50v len:%v\n", "BasePveStage load finish! ", len(result_array))
	return nil
}
