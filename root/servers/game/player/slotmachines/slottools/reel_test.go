package slottools

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	rand.Seed(time.Now().Unix())

	conf := [][]int32{
		{10, 20, 40, 60, 40, 50, 70, 90, 10, 30}, // 列的配置表
		{10, 20, 40, 60, 40, 50, 70, 90, 10, 30},
		{10, 20, 40, 60, 40, 50, 70, 90, 10, 30},
		{10, 20, 40, 60, 40, 50, 70, 90, 10, 30},
		{10, 20, 40, 60, 40, 50, 70, 90, 10, 30},
	}
	ret, e := Reels(conf, 3, -1)
	if e != nil {
		return
	}

	lineIndex := ret.LinePatternIndex([]int32{0, 1, 2, 1, 0})

	ret.Arr = []int32{
		70, 70, 90,
		40, 71, 70,
		60, 40, 10,
		70, 40, 10,
		20, 90, 10,
	}
	fmt.Printf("ret:%v", ret)
	fmt.Printf("line:%v\n", lineIndex)
	Id, c := ret.Continuous(lineIndex, 10, 3)
	fmt.Printf("Id:%v c:%v\n", Id, c)

	ret.RangeRow(func(row, col int, val int32) bool {
		fmt.Printf("row:%v col:%v val:%v \n", row, col, val)
		return false
	})
}
