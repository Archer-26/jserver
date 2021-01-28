package common

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func TestNameMatch(t *testing.T) {
	r, _ := regexp.Compile("Game([0-9]+)_Actor")
	b := r.MatchString("Game232Actor")
	fmt.Print(b)
}

func TestRandWeights(t *testing.T) {
	rand.Seed(time.Now().Unix())
	rands := map[int32]int32{
		1: 10,
		2: 20,
		3: 50,
		4: 20,
	}
	countmap := map[int32]int{}
	for i := 0; i < 10000; i++ {
		v := RandomWeight32(rands)
		countmap[v]++
	}
	fmt.Printf("countmap：%v", countmap)
}
