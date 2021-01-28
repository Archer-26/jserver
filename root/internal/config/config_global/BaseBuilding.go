package config_global

func GetBuildingLevelKey(buildingId int32, buildingLevel int32) int32 {
	return buildingId*1000 + buildingLevel
}

func GetBuildingStarKey(buildingId int32, buildingStar int32) int32 {
	return buildingId*1000 + buildingStar
}

func GetBuildingPropKey(buildingId int32, propId int32, propLevel int32) int32 {
	return buildingId*10000 + propId*1000 + propLevel
}
