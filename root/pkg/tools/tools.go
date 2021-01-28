package tools

import (
	"math/rand"
	"root/pkg/log"
)

type CommonRand struct {
	Id     int
	Chance int32
	Extra  int32
}

func CircleRand(commonArrs []CommonRand) (id int, arrIndex int) {
	arr := make([]CommonRand, len(commonArrs))
	copy(arr, commonArrs)
	var sum int32
	for _, v := range arr {
		sum += v.Chance
	}
	var base int32
	if sum > 0 {
		chance := rand.Intn(int(sum + 1))
		for index, v := range arr {
			if v.Chance == 0 {
				continue
			}
			if int(base+v.Chance) >= chance {
				return v.Id, index
			} else {
				base += v.Chance
			}
		}
	}
	log.KV("arr", arr).ErrorStack(3, "CircleRand Wrong")
	return 0, -1
}
