package slottools

import (
	"fmt"
	"math/rand"
	"root/pkg/log"
)

/*
n行*m列 的轮盘，依次从上至下、从左至右排列成数组 例:
3行*5列 Arr :=[500,600,700, 500,701,100, 200,500,600, 900,702,200, 500,600,100]
	 										╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍
	 										╏ 500 ╎ 500 ╎ 200 ╎ 900 ╎ 500 ╏
	 										╏ 600 ╎ 701 ╎ 500 ╎ 702 ╎ 600 ╏
	 										╏ 700 ╎ 100 ╎ 600 ╎ 200 ╎ 100 ╏
	 										╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍
 图案ID > 10 个位是子类型，十位以上是主类型，例如:71 【7代表主类型、1 代表子类型】，用主类型判断连续
*/
const MAST_TYPE_SCALE = 100

type (
	ReelsArray struct {
		Arr           []int64 // 随机出来的n*m的轮盘
		RandConfIndex []int64 // 每列配置表随机的位置
		ColSize       int64   // 每一列的个数
	}
)

/*
r 每一列的配置
colSize 每一列需要的个数
spec  指定每一列都取固定的一行 -1 代表随机
*/
func Reels(r [][]int64, colSize int64, spec int) (*ReelsArray, error) {
	arr := []int64{}
	randpos := []int64{}
	for _, v := range r {
		col, rand, err := Reel(v, colSize, spec) // 轮询获得每一列的数据
		if err != nil {
			return nil, err
		}
		randpos = append(randpos, rand)
		arr = append(arr, col...)
	}
	return &ReelsArray{Arr: arr, RandConfIndex: randpos, ColSize: colSize}, nil
}

// 得到每条线的图案在Reels中的索引
// line  line的长度表示当前轮盘的列数,从左往右依次代表每列格子索引0,1,2,1,0
func (this *ReelsArray) LinePatternIndex(line []int64) []int {
	// reels的长度必须是line长度的整数倍
	if len(this.Arr)%len(line) != 0 {
		log.KVs(log.Fields{"len(reels)": len(this.Arr), "len(line)": len(line)}).ErrorStack(4, "len(this.Arr)%len(line) != 0")
		panic(nil)
	}
	if len(this.Arr)%int(this.ColSize) != 0 {
		log.KVs(log.Fields{"len(reels)": len(this.Arr), "this.ColSize": this.ColSize}).ErrorStack(4, "len(this.Arr)%int(this.ColSize) != 0")
		panic(nil)
	}
	ret := []int{}
	for k, v := range line {
		if int(v) >= int(this.ColSize) {
			log.KVs(log.Fields{"v": v, "this.ColSize": this.ColSize}).ErrorStack(4, "line val>= rowSize")
			panic(nil)
		}
		i := k*int(this.ColSize) + int(v)
		if i >= len(this.Arr) {
			log.KVs(log.Fields{"i": i, "this.ColSize": this.ColSize, "this.Arr": this.Arr, "line": line}).ErrorStack(4, "reels out of range")
			return nil
		}
		ret = append(ret, i)
	}
	return ret
}

/*
获得一条线中，连续的图案ID和连续数
lineIndex 连线的索引
wildId 万能Id
mincount 最小连续性
id 相连的图案ID，-1表示没有任何图案相连
arr 图案组
*/
func (this *ReelsArray) Continuous(lineIndex []int, wildId int64, mincount int64) (id int64, arr []int64) {
	if this.Arr == nil || len(this.Arr) == 0 || len(lineIndex) == 0 {
		log.KVs(log.Fields{"ReelsArray": this.String(), "lineIndex": lineIndex}).ErrorStack(4, "data error ")
		return -1, nil
	}

	currentVal := this.Arr[lineIndex[0]] / MAST_TYPE_SCALE
	id = currentVal

	if currentVal == wildId {
		log.KVs(log.Fields{"currentVal": currentVal, "wildId": wildId}).ErrorStack(4, "arr's frist val == wild ")
		return -1, nil
	}
	arr = make([]int64, 0)
	for i := 0; i < len(lineIndex); i++ {
		if lineIndex[i] >= len(this.Arr) {
			log.KVs(log.Fields{"lineIndex": lineIndex, "this.Arr": this.Arr}).ErrorStack(4, "out of range ")
			return
		}
		t := this.Arr[lineIndex[i]] / MAST_TYPE_SCALE
		if currentVal == t || wildId == this.Arr[lineIndex[i]] {
			arr = append(arr, this.Arr[lineIndex[i]])
		} else {
			break
		}
	}
	if len(arr) >= int(mincount) {
		return this.Arr[lineIndex[0]], arr
	} else {
		return -1, nil
	}
}

func (this *ReelsArray) RangeRow(f func(row, col int, val int64) bool) {
	if len(this.Arr)%int(this.ColSize) != 0 {
		//log.KVs(log.Fields{"this.Arr": this.Arr, "this.ColSize": this.ColSize}).ErrorStack(3, "reels n*m size error")
		return
	}
	col_len := len(this.Arr) / int(this.ColSize)
	for row := 0; row < int(this.ColSize); row++ {
		for col := 0; col < col_len; col++ {
			index := row + (col * int(this.ColSize))
			if f(row, col, this.Arr[index]) {
				return
			}
		}
	}
}

// 输入一列轮盘，随机获得需要的结果
// r 单列轮盘配置
// c 需要输出几个
// spec 指定停留在某一列
func Reel(r []int64, c int64, spec int) ([]int64, int64, error) {
	ret := []int64{}
	l := len(r)
	if l <= int(c) {
		return nil, -1, fmt.Errorf("len(wheel) < c, len:%v c:%v", l, c)
	}
	randPos := spec
	if randPos == -1 {
		randPos = rand.Intn(l)
	}

	for i := 0; i < int(c); i++ {
		if randPos >= l {
			randPos = 0
		}
		ret = append(ret, r[randPos])
		randPos++
	}
	return ret, int64(randPos), nil
}

func (this *ReelsArray) String() string {
	rowstr := "\n"
	for row := 0; row < int(this.ColSize); row++ {
		for i := row; i < len(this.Arr); i += int(this.ColSize) {
			rowstr += fmt.Sprintf("%-4v ", this.Arr[i])
		}
		rowstr += "\n"
	}
	rowstr += fmt.Sprintf("%v", this.RandConfIndex)
	return rowstr
}

// index周围4个位置的索引[上,下,左,右]
func (this *ReelsArray) AroundPos(index int) []int {
	up := index - 1
	down := index + 1

	row := index % int(this.ColSize)
	if row == 0 {
		up = -1
	}
	if row == 2 {
		down = -1
	}
	return []int{up, down, index - int(this.ColSize), index + int(this.ColSize)}
}

// 位置是否有效
func (this *ReelsArray) PosCheck(pos int) bool {
	if 0 <= pos && pos < len(this.Arr) {
		return true
	} else {
		return false
	}
}
