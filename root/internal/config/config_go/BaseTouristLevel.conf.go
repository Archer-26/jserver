package config_go

// Code generated by excelExoprt. DO NOT EDIT.
// source: BaseTouristLevel.xlsx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// array and map
var _BaseTouristLevelMap = map[int64]*BaseTouristLevel{}
var _BaseTouristLevelArray = []*BaseTouristLevel{}

type BaseTouristLevel struct {
	data *_BaseTouristLevel
}

// 类型结构
type _BaseTouristLevel struct {
	INT_ID             int64     // ID
	INT_TOURIST_ID     int64     // 游客类型
	INT_LEVEL          int64     // 等级
	INT_EXPECTED       int64     // 预期
	INT_MAX_HAPPINESS  int64     // 欢乐值上限
	INT_ENERGY         int64     // 体力
	INT_RUN_SPEED      int64     // 行进速度
	INT_BUILDING_LOSS  int64     // 对建筑损耗
	INT_SERVER_REQUIRE int64     // 服务值
	INT_PRAISE         int64     // 称赞值
	INT_CRITIC         int64     // 批评值
	INT_HERO_LOSS      int64     // 对主管体力需求
	ARRAYINT_SKILL_ID  array_int // 技能ID
	INT_VICTORY_POINT  int64     // 提供分数
}

func (c *BaseTouristLevel) ID() int64             { return c.data.INT_ID }             //ID
func (c *BaseTouristLevel) TOURIST_ID() int64     { return c.data.INT_TOURIST_ID }     //游客类型
func (c *BaseTouristLevel) LEVEL() int64          { return c.data.INT_LEVEL }          //等级
func (c *BaseTouristLevel) EXPECTED() int64       { return c.data.INT_EXPECTED }       //预期
func (c *BaseTouristLevel) MAX_HAPPINESS() int64  { return c.data.INT_MAX_HAPPINESS }  //欢乐值上限
func (c *BaseTouristLevel) ENERGY() int64         { return c.data.INT_ENERGY }         //体力
func (c *BaseTouristLevel) RUN_SPEED() int64      { return c.data.INT_RUN_SPEED }      //行进速度
func (c *BaseTouristLevel) BUILDING_LOSS() int64  { return c.data.INT_BUILDING_LOSS }  //对建筑损耗
func (c *BaseTouristLevel) SERVER_REQUIRE() int64 { return c.data.INT_SERVER_REQUIRE } //服务值
func (c *BaseTouristLevel) PRAISE() int64         { return c.data.INT_PRAISE }         //称赞值
func (c *BaseTouristLevel) CRITIC() int64         { return c.data.INT_CRITIC }         //批评值
func (c *BaseTouristLevel) HERO_LOSS() int64      { return c.data.INT_HERO_LOSS }      //对主管体力需求
//技能ID
func (c *BaseTouristLevel) Len_SKILL_ID() int          { return c.data.ARRAYINT_SKILL_ID.Len() }
func (c *BaseTouristLevel) Get_SKILL_ID(key int) int64 { return c.data.ARRAYINT_SKILL_ID.Get(key) }
func (c *BaseTouristLevel) Range_SKILL_ID(fn func(int, int64) (stop bool)) {
	c.data.ARRAYINT_SKILL_ID.Range(fn)
}
func (c *BaseTouristLevel) Copy_SKILL_ID() array_int { return c.data.ARRAYINT_SKILL_ID.Copy() }
func (c *BaseTouristLevel) VICTORY_POINT() int64     { return c.data.INT_VICTORY_POINT } //提供分数

func HasBaseTouristLevel(key int64) bool {
	_, ok := _BaseTouristLevelMap[key]
	return ok
}

func GetBaseTouristLevel(key int64) *BaseTouristLevel {
	return _BaseTouristLevelMap[key]
}

func RangeBaseTouristLevel(fn func(i int, row *BaseTouristLevel) (stop bool)) {
	for i, row := range _BaseTouristLevelArray {
		if fn(i, row) {
			break
		}
	}
}

func LenBaseTouristLevel() int { return len(_BaseTouristLevelArray) }

func init() {
	loadfn["BaseTouristLevel"] = loadBaseTouristLevel
}

func loadBaseTouristLevel(dir string) error {
	data, err := ioutil.ReadFile(path.Join(dir, "BaseTouristLevel.json"))
	if err != nil {
		return fmt.Errorf("file=%v read err=%v", err.Error())
	}

	datas := []*_BaseTouristLevel{}
	err = json.Unmarshal(data, &datas)
	if err != nil {
		return fmt.Errorf("file=%v parse err=%v", err.Error())
	}

	result_array := []*BaseTouristLevel{}
	result_map := map[int64]*BaseTouristLevel{}
	for _, row := range datas {
		data := &BaseTouristLevel{data: row}
		result_array = append(result_array, data)
		result_map[row.INT_ID] = data
	}
	_BaseTouristLevelArray = result_array
	_BaseTouristLevelMap = result_map
	fmt.Printf("%-50v len:%v\n", "BaseTouristLevel load finish! ", len(result_array))
	return nil
}
