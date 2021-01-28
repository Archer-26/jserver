package items

import (
	"root/internal/common"
	"root/internal/config/config_global"
	"root/internal/config/config_go"
	"root/pkg/log"
	"root/servers/game/interfaces"
)

type Building struct {
	BaseItem
}

func newBuilding() *Building {
	return &Building{}
}

func (this *Building) Add(role interfaces.IPlayer, itemId int32, itemCount int64, itemLog *common.ItemLog) {
	baseItem := config_go.GetBaseItem(int64(itemId))
	starKey := config_global.GetBuildingStarKey(int32(baseItem.BUILDING()), 1)
	buildingStar := config_go.GetBaseBuildingStar(int64(starKey))
	if buildingStar == nil {
		log.KVs(log.Fields{"rid": role.RID(), "starKey": starKey}).Error("buildingStar no this key")
		return
	}
	for i := 1; i <= int(itemCount); i++ {
		role.Building().AddBuilding(int32(baseItem.BUILDING()), 1, 1, itemLog)
	}
}
