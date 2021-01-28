package config_global

import (
	"root/internal/config/config_go"
)

var BuildingWithProp map[int64]map[int64]bool // buildingId => map[ATTRId]
func init() {
	config_go.RegisterPackageFn(config_go.ExtLoad_PackageFn, Init)
}
func Init() {
	packageBuildingProp()
}

func packageBuildingProp() {
	var data = map[int64]map[int64]bool{}
	config_go.RangeBaseBuildingAttr(func(i int, row *config_go.BaseBuildingAttr) bool {
		props, ok := data[row.BUILDING_ID()]
		if !ok {
			data[row.BUILDING_ID()] = map[int64]bool{}
			props = data[row.BUILDING_ID()]
		}
		props[row.ATTR()] = true
		return false
	})
	BuildingWithProp = data
}
