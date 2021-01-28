package slotmachines

import (
	"fmt"
	"root/servers/game/player/slotmachines/slottools"
	"testing"
)

func TestName(t *testing.T) {
	reels := &slottools.ReelsArray{}
	reels.Arr = []int32{
		301, 500, 300,
		302, 700, 300,
		305, 400, 300,
		400, 300, 300,
		400, 300, 300,
	}
	reels.ColSize = 3
	loginfo := fmt.Sprintf("bounus 特殊1 test:%v", reels.String())
	fmt.Printf("%v\n", loginfo)
	CheckExtra3Linked(reels)
	loginfo = fmt.Sprintf("bounus 特殊2 test:%v", reels.String())
	fmt.Printf("%v\n", loginfo)
}
