package common

type REASON int32

const (
	SLOT_BET          REASON = 1
	CMD               REASON = 2
	USE_ITEM          REASON = 3
	BUILDING_LEVEL_UP REASON = 4
	BUILDING_START_UP REASON = 5
	ROLE_LEVEL_UP     REASON = 6
	CHECKPOINT_FINISH REASON = 7
	DICE              REASON = 8
	BUILDING_PROP_UP  REASON = 9
)

type ItemLog struct {
	Id           uint64
	ModifyReason REASON
}

func NewItemLog(reason REASON) *ItemLog {
	return &ItemLog{
		Id:           GetNewId(),
		ModifyReason: reason,
	}
}

func NewItemLogWithId(id uint64, reason REASON) *ItemLog {
	return &ItemLog{
		Id:           id,
		ModifyReason: reason,
	}
}
