package config_global

import (
	"root/internal/config/config_go"
	"root/pkg/log"
	"root/pkg/tools"
)

func GetRewardByPool(poolId int64, roleLevel int64, itemMap map[int32]int64) {
	pool := config_go.GetBaseRewardPool(poolId)
	if pool == nil {
		log.KVs(log.Fields{"poolId": poolId}).ErrorStack(3, "no this pool")
		return
	}

	for i := 0; i < pool.Len_REWARD_ID(); i++ {
		if int64(roleLevel) < pool.Get_MIN_LV(i) || int64(roleLevel) > pool.Get_MAX_LV(i) {
			continue
		}
		GetRewardById(int32(pool.Get_REWARD_ID(i)), itemMap)
	}
	return
}

func GetRewardById(rewardId int32, itemMap map[int32]int64) {
	reward := config_go.GetBaseReward(int64(rewardId))
	if reward == nil {
		log.KVs(log.Fields{"rewardId": rewardId}).ErrorStack(3, "no this rewardId")
		return
	}

	switch reward.PROB_TYPE() {
	case 1: //所有数据与概率比较
		for i := 0; i < reward.Len_ITEM_ID(); i++ {
			if tools.Probability10000(int(reward.Get_PROB(i))) {
				itemMap[int32(reward.Get_ITEM_ID(i))] += reward.Get_ITEM_COUNT(i)
			}
		}
	case 2: // 遍历所有
		arr := []tools.CommonRand{}
		for i := 0; i < reward.Len_ITEM_ID(); i++ {
			arr = append(arr, tools.CommonRand{Id: i, Chance: int32(reward.Get_PROB(i))})
		}
		index, _ := tools.CircleRand(arr)
		itemMap[int32(reward.Get_ITEM_ID(index))] += reward.Get_ITEM_COUNT(index)
	}
}
