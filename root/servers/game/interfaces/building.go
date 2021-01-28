package interfaces

import (
	"root/internal/common"
	"root/servers/internal/models"
)

type IBuilding interface {
	GetBuilding(buildingId int32) *models.Building
	AddBuilding(buildingId int32, level int32, star int32, itemLog *common.ItemLog)
	SetBuildingLevel(buildingId int32, newLevel int32, itemLog *common.ItemLog)
	SetBuildingStar(buildingId int32, newStar int32, itemLog *common.ItemLog)
	SetBuildingPropLevel(buildingId int32, propId int32, newPropLevel int32, itemLog *common.ItemLog)
}
