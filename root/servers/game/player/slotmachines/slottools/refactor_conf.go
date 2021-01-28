package slottools

import (
	"root/internal/common"
	"root/internal/config/config_go"
)

type Slot_Base_Reel struct {
	Solt_Base_Reels  [][]int64
	Solt_Bonus_Index [][]int64 // 每一列bonus所在位置
}

const (
	BonusType = 3
	BaseBonus = 300
	Scatter   = 200 // Scatter

	JackpotMINI  = 321
	JackpotMINOR = 322
	JackpotMAJOR = 323
	JackpotGRAND = 324

	JackpotAlter = 310 // jackpot改变后的值

	WildID     = 100
	SevenType  = 9 // 7相关图案的类型
	Special777 = 903
)

var (
	///////////////////////////////////////////// 老虎机相关 ////////////////////////////////////////////////
	SlotBaseReel   *Slot_Base_Reel
	SlotSelectReel *Slot_Base_Reel
	SlotFreeReel   *Slot_Base_Reel

	SlotProbabilityRoole3 map[int64]int64
	SlotProbabilityRoole  map[int64]int64

	SlotOption1 [][]int64
	SlotOption2 [][]int64
	SlotReel1_5 map[int64][][]int64 // [respin次数][reel1,reel2,reel3,reel4,reel5][需要数量，是否邻近]
)

func init() {
	common.RegistRefactorFun(SLot)
}
func SLot() {
	tempBaseRell := &Slot_Base_Reel{
		Solt_Base_Reels:  make([][]int64, 0),
		Solt_Bonus_Index: make([][]int64, 0),
	}
	tempSelectRell := &Slot_Base_Reel{
		Solt_Base_Reels:  make([][]int64, 0),
		Solt_Bonus_Index: make([][]int64, 0),
	}
	tempFreeRell := &Slot_Base_Reel{
		Solt_Base_Reels:  make([][]int64, 0),
		Solt_Bonus_Index: make([][]int64, 0),
	}
	tempOption1 := make([][]int64, 0)
	tempOption2 := make([][]int64, 0)
	config_go.RangeSlotCherryReels(func(i int, row *config_go.SlotCherryReels) bool {
		// 处理基础表盘
		reel := row.Copy_REEL()
		tempBaseRell.Solt_Base_Reels = append(tempBaseRell.Solt_Base_Reels, reel)
		tempBaseRell.Solt_Bonus_Index = append(tempBaseRell.Solt_Bonus_Index, []int64{})
		for index, v := range reel {
			if v == BonusType*100 {
				tempBaseRell.Solt_Bonus_Index[i] = append(tempBaseRell.Solt_Bonus_Index[i], int64(index))
			}
		}

		// 处理option表盘
		tempOption1 = append(tempOption1, row.Copy_OPTION_1())
		tempOption2 = append(tempOption2, row.Copy_OPTION_1())

		// select 表盘
		selection := row.Copy_SLELCTION()
		tempSelectRell.Solt_Base_Reels = append(tempSelectRell.Solt_Base_Reels, selection)
		tempSelectRell.Solt_Bonus_Index = append(tempSelectRell.Solt_Bonus_Index, []int64{})
		for index, v := range selection {
			if v == BonusType*100 {
				tempSelectRell.Solt_Bonus_Index[i] = append(tempSelectRell.Solt_Bonus_Index[i], int64(index))
			}
		}
		// free 表盘
		free := row.Copy_SLELCTION()
		tempFreeRell.Solt_Base_Reels = append(tempFreeRell.Solt_Base_Reels, free)
		tempFreeRell.Solt_Bonus_Index = append(tempFreeRell.Solt_Bonus_Index, []int64{})
		for index, v := range free {
			if v == BonusType*100 {
				tempFreeRell.Solt_Bonus_Index[i] = append(tempFreeRell.Solt_Bonus_Index[i], int64(index))
			}
		}
		return false
	})
	SlotBaseReel = tempBaseRell
	SlotSelectReel = tempSelectRell
	SlotFreeReel = tempFreeRell
	SlotOption1 = tempOption1
	SlotOption2 = tempOption2

	tempPattern3 := map[int64]int64{}
	tempPattern := map[int64]int64{}
	config_go.RangeSlotCherryPattern(func(i int, row *config_go.SlotCherryPattern) bool {
		if row.ID()/100 == BonusType {
			tempPattern3[row.ID()] = row.PROBABILITY_REEL3()
			tempPattern[row.ID()] = row.PROBABILITY_OTHER_REEL()
		}
		return false
	})
	SlotProbabilityRoole3 = tempPattern3
	SlotProbabilityRoole = tempPattern

	reel1_5 := map[int64][][]int64{}
	config_go.RangeCherryRespin(func(i int, row *config_go.CherryRespin) bool {
		reel1_5[row.ID()] = [][]int64{}
		reel1_5[row.ID()] = append(reel1_5[row.ID()], row.Copy_REEL1())
		reel1_5[row.ID()] = append(reel1_5[row.ID()], row.Copy_REEL2())
		reel1_5[row.ID()] = append(reel1_5[row.ID()], row.Copy_REEL3())
		reel1_5[row.ID()] = append(reel1_5[row.ID()], row.Copy_REEL4())
		reel1_5[row.ID()] = append(reel1_5[row.ID()], row.Copy_REEL5())
		return false
	})
	SlotReel1_5 = reel1_5
	//test()
	//test2()
}