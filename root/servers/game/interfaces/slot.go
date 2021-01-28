package interfaces

import "root/servers/internal/message"

type ISlot interface {
	Bet(betVal int64, lv int64, triggerType int, DiscountRate int, BaseMsg *message.Base, BonusMsg *message.Bonus) (gold int64, base, bonus, scatter map[int32]int64)
	Select1(triggerType int, BonusMsg *message.Bonus) (randrow int, BonusOdds int64)
	SetLastBetLv(lv int64)
	LastBetLv() int64
	Select() bool
	Save()
}
