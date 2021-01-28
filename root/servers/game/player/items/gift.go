package items

import (
	"root/internal/common"
	"root/internal/config/config_global"
	"root/internal/config/config_go"
	"root/servers/game/interfaces"
	//"root/config/config_global"
	//"root/config/config_go"
	"root/servers/internal/message"
)

type Gift struct {
	Normal
}

func newGift() *Gift {
	return &Gift{}
}

func (this *Gift) Use(role interfaces.IPlayer, itemId int32, itemCount int64, extParams *message.UseExtParams, itemLog *common.ItemLog) (flag bool) {
	baseItem := config_go.GetBaseItem(int64(itemId))
	itemMap := make(map[int32]int64)
	for i := 1; i <= int(itemCount); i++ {
		config_global.GetRewardByPool(baseItem.REWARD_POOL(), int64(role.Lv()), itemMap)
	}
	role.AddItemMap(itemMap, itemLog)
	return true
}
