package config_go

// Code generated by excelExoprt. DO NOT EDIT.
// source: FunnyGameLevelDetail.xlsx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// array and map
var _FunnyGameLevelDetailMap = map[int64]*FunnyGameLevelDetail{}
var _FunnyGameLevelDetailArray = []*FunnyGameLevelDetail{}

type FunnyGameLevelDetail struct {
	data *_FunnyGameLevelDetail
}

// 类型结构
type _FunnyGameLevelDetail struct {
	INT_ID                  int64     // ID
	INT_LEVEL_INDEX         int64     // 关卡索引ID
	ARRAYINT_INDEX_POSITION array_int // 索引位置
	INT_ROAD_TYPE           int64     // 道路类型
	INT_PROP_TYPE           int64     // 道具类型
	INT_PROP_COUNT          int64     // 道具的数量
	INT_PARAMETER           int64     // 道具参数
}

func (c *FunnyGameLevelDetail) ID() int64          { return c.data.INT_ID }          //ID
func (c *FunnyGameLevelDetail) LEVEL_INDEX() int64 { return c.data.INT_LEVEL_INDEX } //关卡索引ID
//索引位置
func (c *FunnyGameLevelDetail) Len_INDEX_POSITION() int { return c.data.ARRAYINT_INDEX_POSITION.Len() }
func (c *FunnyGameLevelDetail) Get_INDEX_POSITION(key int) int64 {
	return c.data.ARRAYINT_INDEX_POSITION.Get(key)
}
func (c *FunnyGameLevelDetail) Range_INDEX_POSITION(fn func(int, int64) (stop bool)) {
	c.data.ARRAYINT_INDEX_POSITION.Range(fn)
}
func (c *FunnyGameLevelDetail) Copy_INDEX_POSITION() array_int {
	return c.data.ARRAYINT_INDEX_POSITION.Copy()
}
func (c *FunnyGameLevelDetail) ROAD_TYPE() int64  { return c.data.INT_ROAD_TYPE }  //道路类型
func (c *FunnyGameLevelDetail) PROP_TYPE() int64  { return c.data.INT_PROP_TYPE }  //道具类型
func (c *FunnyGameLevelDetail) PROP_COUNT() int64 { return c.data.INT_PROP_COUNT } //道具的数量
func (c *FunnyGameLevelDetail) PARAMETER() int64  { return c.data.INT_PARAMETER }  //道具参数

func HasFunnyGameLevelDetail(key int64) bool {
	_, ok := _FunnyGameLevelDetailMap[key]
	return ok
}

func GetFunnyGameLevelDetail(key int64) *FunnyGameLevelDetail {
	return _FunnyGameLevelDetailMap[key]
}

func RangeFunnyGameLevelDetail(fn func(i int, row *FunnyGameLevelDetail) (stop bool)) {
	for i, row := range _FunnyGameLevelDetailArray {
		if fn(i, row) {
			break
		}
	}
}

func LenFunnyGameLevelDetail() int { return len(_FunnyGameLevelDetailArray) }

func init() {
	loadfn["FunnyGameLevelDetail"] = loadFunnyGameLevelDetail
}

func loadFunnyGameLevelDetail(dir string) error {
	data, err := ioutil.ReadFile(path.Join(dir, "FunnyGameLevelDetail.json"))
	if err != nil {
		return fmt.Errorf("file=%v read err=%v", err.Error())
	}

	datas := []*_FunnyGameLevelDetail{}
	err = json.Unmarshal(data, &datas)
	if err != nil {
		return fmt.Errorf("file=%v parse err=%v", err.Error())
	}

	result_array := []*FunnyGameLevelDetail{}
	result_map := map[int64]*FunnyGameLevelDetail{}
	for _, row := range datas {
		data := &FunnyGameLevelDetail{data: row}
		result_array = append(result_array, data)
		result_map[row.INT_ID] = data
	}
	_FunnyGameLevelDetailArray = result_array
	_FunnyGameLevelDetailMap = result_map
	fmt.Printf("%-50v len:%v\n", "FunnyGameLevelDetail load finish! ", len(result_array))
	return nil
}
