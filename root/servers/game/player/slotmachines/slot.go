package slotmachines

import (
	"fmt"
	"math/rand"
	"root/internal/common"
	"root/internal/config/config_global"
	"root/internal/config/config_go"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/pkg/tools"
	"root/servers/game/interfaces"
	"root/servers/game/player/slotmachines/slottools"
	"root/servers/game/player/slotmachines/slottype"
	"root/servers/internal/message"
	"root/servers/internal/models"
	"root/servers/internal/mysql/i_mysql"
	"strconv"
	"strings"
)

type (
	Slot struct {
		RID  int64
		game interfaces.IGame
		*models.ModelSlot
		jackpot map[int64][]int64
	}
)

func (this *Slot) Init(rid int64, g interfaces.IGame) {
	this.RID = rid
	this.game = g
	this.ModelSlot = models.NewModelSlot(rid)
	this.jackpot = make(map[int64][]int64, 0)
	majormin := int64(tools.Randx_y(int(config_go.GetCherryJackPot(slottools.JackpotMAJOR).InitMin()), int(config_go.GetCherryJackPot(slottools.JackpotMAJOR).InitMax())))
	majormax := int64(tools.Randx_y(int(config_go.GetCherryJackPot(slottools.JackpotMAJOR).RstMin()), int(config_go.GetCherryJackPot(slottools.JackpotMAJOR).RstMax())))

	grandmin := int64(tools.Randx_y(int(config_go.GetCherryJackPot(slottools.JackpotGRAND).InitMin()), int(config_go.GetCherryJackPot(slottools.JackpotGRAND).InitMax())))
	grandmax := int64(tools.Randx_y(int(config_go.GetCherryJackPot(slottools.JackpotGRAND).RstMin()), int(config_go.GetCherryJackPot(slottools.JackpotGRAND).RstMax())))

	this.jackpot[slottools.JackpotMAJOR] = []int64{majormin, majormax}
	this.jackpot[slottools.JackpotGRAND] = []int64{grandmin, grandmax}
}

func (this *Slot) Sync() {
	slotpb := &message.NotifySlotInfo{
		LastBetOdds: this.LastBet, // 最后一次押注赔率
	}
	msg := network.NewPbMessage(slotpb, message.MSG_NOTIFY_SLOT_INFO.Int32())
	this.Role().SendMsg(msg)
}

func (this *Slot) SetLastBetLv(lv int64) {
	this.LastBet = int64(lv)
}
func (this *Slot) LastBetLv() int64 { return int64(this.LastBet) }
func (this *Slot) Select() bool     { return this.IsSelect }

// odds 押注赔率
// reward 返回奖励列表
func (this *Slot) Bet(BetVal int64, lv int64, triggerType int, DiscountRate int, BaseMsg *message.Base, BonusMsg *message.Bonus) (gold int64, base, bonus, scatter map[int32]int64) {
	var Conf_Reel *slottools.Slot_Base_Reel
	if triggerType == slottype.Trigger_Respin_Base {
		Conf_Reel = slottools.SlotBaseReel
	} else if triggerType == slottype.Trigger_Respin_Free {
		Conf_Reel = slottools.SlotFreeReel
		this.IsSelect = false
	} else {
		log.KV("triggerType", slottype.Trigger_Respin_Free).ErrorStack(3, "triggerType error")
		return
	}
	// 把轮盘中的Bonus根据配置的权重替换
	for i, indexs := range Conf_Reel.Solt_Bonus_Index {
		for _, index := range indexs {
			if i == 2 { // 第三列单独用一套配置
				Conf_Reel.Solt_Base_Reels[i][int(index)] = common.RandomWeight32(slottools.SlotProbabilityRoole3)
			} else {
				Conf_Reel.Solt_Base_Reels[i][int(index)] = common.RandomWeight32(slottools.SlotProbabilityRoole)
			}
		}
	}

	// 获得基础轮盘
	reels, e := slottools.Reels(Conf_Reel.Solt_Base_Reels, 3, -1)
	if e != nil {
		log.KV("error", e).ErrorStack(3, "slottools.Reels error ")
		return
	}

	//if triggerType == slottype.Trigger_Respin_Base {
	//	reels.Arr = []int64{
	//		301, 500, 900,
	//		302, 200, 700,
	//		200, 400, 600,
	//		500, 400, 200,
	//		500, 900, 500,
	//	}
	//}

	BaseMsg.BaseReels = reels.Arr
	BaseMsg.RandIndex = reels.RandConfIndex
	loginfo := fmt.Sprintf("%v", reels)
	log.Debug(colorized.Yellow(loginfo))

	RewardOddsMap := map[int64]int64{} // [中奖的线Id]Odds
	totalGold := int64(0)
	totalOdds := int64(0)
	baseItemMap := map[int32]int64{}
	bonusItemMap := make(map[int32]int64, 0)
	scatterItemMap := make(map[int32]int64, 0)
	BaseMsg.RewardLines = &message.RewardLine{Alter: make(map[int64]int64)}

	// 依次判断每条线是否中奖 ////////////////////////////////////////////////////////////////////////////////////////////
	config_go.RangeSlotCherryLine(func(lineId int, row *config_go.SlotCherryLine) bool {
		line := row.Copy_Line()
		line_index := reels.LinePatternIndex(line) // 获得线的索引
		partternId, arr := reels.Continuous(line_index, slottools.WildID, 3)
		if partternId == -1 {
			return false
		}

		// 7相关图案特殊判断
		if patternType(partternId) == slottools.SevenType {
			if len(arr) < 3 {
				log.KV("len(arr)", len(arr)).ErrorStack(3, "len(arr) < 3 error ")
				return false
			}
			sevenId := arr[0]
			for _, v := range arr {
				if sevenId != v { // 只要有一个7不一样，就以777图案获得奖励
					partternId = slottools.Special777
					break
				}
			}
		}

		// 统计奖励
		rewardOdds := config_go.GetSlotCherryPattern(partternId).Copy_LINE()
		if rewardOdds == nil {
			log.KVs(log.Fields{"partternId": partternId, "arr": arr, "line": line}).Warn("partternId not reward")
			return false
		}
		if len(arr) > len(rewardOdds) {
			log.KVs(log.Fields{"partternId": partternId, "arr": arr, "rewardOdds": rewardOdds, "line": line}).ErrorStack(3, "len(arr) >= len(rewardOdds)")
			return false
		}
		BaseMsg.RewardLines.Alter[int64(lineId)] = rewardOdds[len(arr)-1]
		RewardOddsMap[int64(lineId)] = rewardOdds[len(arr)-1]
		totalOdds += int64(rewardOdds[len(arr)-1])
		totalGold += int64(rewardOdds[len(arr)-1]) * BetVal / int64(config_go.LenSlotCherryLine())

		// 连线触发道具掉落
		produceId := 1000 + int64(len(arr))
		itemMap(lv, produceId, DiscountRate, baseItemMap)
		log.KVs(log.Fields{"partternId": partternId, "arr": arr, "line": line}).Info(colorized.Yellow("line bingo!"))
		return false
	})

	BaseMsg.LineGlod = totalGold
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// 先检查Scatter Scatter没中，再判断Bonus
	this.IsSelect = ScatterPlay(reels)
	// respin玩法检测
	if !this.IsSelect {
		combineId, rowId, BonusOdds, jackpotOdds := this.BonusPlay(reels, triggerType, BonusMsg)
		log.KVs(log.Fields{"combineId": combineId, "rowId": rowId, "BonusOdds": BonusOdds, "jackpotOdds": jackpotOdds}).Info("bonus")
		bonusGold := (BonusOdds + jackpotOdds) * BetVal
		totalGold += bonusGold
		totalOdds += BonusOdds + jackpotOdds
		BonusMsg.BonusGlod = bonusGold
		if combineId != "" {
			// bonus触发道具掉落
			produceId := int64(2000 + triggerType)
			itemMap(lv, produceId, DiscountRate, bonusItemMap)
		}
	} else {
		// scatter触发道具掉落
		if triggerType == slottype.Trigger_Respin_Base {
			produceId := int64(3000)
			itemMap(lv, produceId, DiscountRate, scatterItemMap)
		}
	}
	log.KVs(log.Fields{"IsSelect": this.IsSelect}).Info("ScatterPlay")
	// 检测大奖昵称
	config_go.RangeCherryWinType(func(i int, row *config_go.CherryWinType) bool {
		if int64(row.MIN()) <= totalOdds && totalOdds <= int64(row.MAX()) {
			BaseMsg.WinTypeId = row.ID()
			itemMap(lv, 4000+BaseMsg.WinTypeId, DiscountRate, baseItemMap)
			return true
		}
		return false
	})

	BaseMsg.Items = baseItemMap
	BonusMsg.Items = bonusItemMap
	this.LastBet = BetVal
	return totalGold, baseItemMap, bonusItemMap, scatterItemMap
}

func itemMap(lv, produceId int64, dis int, m map[int32]int64) {
	poolId := config_go.GetCherryProduce(produceId).RewardPoolId()
	tempMap := map[int32]int64{}
	config_global.GetRewardByPool(poolId, lv, tempMap)
	discount := config_go.GetCherryProduce(produceId).Discount()
	if discount == 0 {
		log.KVs(log.Fields{"produceId": produceId}).ErrorStack(3, "discount == 0")
		discount = 1
	}
	baseCount := dis / int(discount)
	for k, _ := range tempMap {
		m[k] += int64(baseCount)
	}
}

// 选择Bonus玩法
func (this *Slot) Select1(triggerType int, BonusMsg *message.Bonus) (randrow int, Odds int64) {
	this.IsSelect = false
	// 把轮盘中的Bonus根据配置的权重替换
	for i, indexs := range slottools.SlotSelectReel.Solt_Bonus_Index {
		for _, index := range indexs {
			if i == 2 { // 第三列单独用一套配置
				slottools.SlotSelectReel.Solt_Base_Reels[i][int(index)] = common.RandomWeight32(slottools.SlotProbabilityRoole3)
			} else {
				slottools.SlotSelectReel.Solt_Base_Reels[i][int(index)] = common.RandomWeight32(slottools.SlotProbabilityRoole)
			}
		}
	}

	// 获得Selection轮盘
	randrow = rand.Intn(len(slottools.SlotSelectReel.Solt_Base_Reels[0]))
	reels, e := slottools.Reels(slottools.SlotSelectReel.Solt_Base_Reels, 3, randrow)
	if e != nil {
		log.KV("error", e).ErrorStack(3, "slottools.Reels error ")
		return 0, 0
	}

	loginfo := fmt.Sprintf("%v", reels)
	log.Debug(colorized.Yellow(loginfo))
	combineId, rowId, BonusOdds, jackpotOdds := this.BonusPlay(reels, triggerType, BonusMsg)
	log.KVs(log.Fields{"combineId": combineId, "rowId": rowId, "BonusOdds": BonusOdds, "jackpotOdds": jackpotOdds}).Info("respin")
	return randrow, BonusOdds + jackpotOdds
}

// Scatter触发检测
func ScatterPlay(reels *slottools.ReelsArray) bool {
	count := 0
	for i := 3; i <= 11; i++ {
		if reels.Arr[i] == slottools.Scatter {
			count++
		}
	}
	if count < 3 {
		return false
	} else {
		return true
	}
}

// BonusPlay触发检查
func (this *Slot) BonusPlay(reels *slottools.ReelsArray, triggerType int, BonusMsg *message.Bonus) (combineId string, rowId int64, BonusOdds int64, jackpotOdds int64) {
	// 检查每一行Bonus的数量,配置上避免同时多行出现bonus>=3个的情况
	combineId = ""
	rowId = -1
	jackpot := map[int64]int{}
	for row := 0; row < int(reels.ColSize); row++ {
		tempId := ""
		for i := 0; i < len(reels.Arr)/int(reels.ColSize); i++ {
			index := row + i*int(reels.ColSize)
			if patternType(reels.Arr[index]) == slottools.BonusType {
				tempId += strconv.Itoa(i + 1)
			}
		}
		if len(tempId) >= 3 { // 如果某一行满足条件直接break，就不需要再找了
			combineId = tempId
			rowId = int64(row)
			break
		}
	}

	if combineId == "" {
		// 正常情况，没有触发respin
		return "", -1, 0, 0
	}

	BonusMsg.CombineId = combineId
	BonusMsg.TriggerRow = rowId
	// 合并相关 /////////////////////////////////////////////////////////////////////////////////////////////////////////
	combine_conf := config_go.GetCherryBonusCombine(combineId)
	if combine_conf == nil {
		log.KV("combineId", combineId).ErrorStack(3, "not found GetCherryBonusCombine")
		return "", -1, 0, 0
	}
	id_array := strings.Split(combineId, "")
	combine_array := combine_conf.Copy_Combine()
	gameCount := len(id_array)
	if len(id_array) != combine_array.Len() {
		log.KVs(log.Fields{"combineId": combineId, "combine_array": combine_array}).ErrorStack(3, "len(id_array) != combine_array.Len()")
		return "", -1, 0, 0
	}

	tempflags := make([]int64, len(reels.Arr), len(reels.Arr)) // 合并后复制的一张表，把无效bonus排除，留下有效的bonus
	// 检查是否需要合并
	for i := 0; i < len(id_array); i++ {
		idpos, _ := strconv.Atoi(id_array[i])
		combpos, _ := strconv.Atoi(combine_array[i])
		real_idpos := int(rowId) + int(reels.ColSize)*(idpos-1)
		real_combpos := int(rowId) + int(reels.ColSize)*(combpos-1)
		if idpos != combpos {
			// 交换位置
			reels.Arr[real_idpos], reels.Arr[real_combpos] = reels.Arr[real_combpos], reels.Arr[real_idpos]
		}
		tempflags[real_combpos] = reels.Arr[real_combpos]
		if pattern := reels.Arr[real_combpos]; pattern >= slottools.JackpotMINI {
			jackpot[pattern] = jackpot[pattern] + 1
			reels.Arr[real_combpos] = slottools.JackpotAlter
		}
	}

	loginfo := fmt.Sprintf("bounus combine:%v", reels)
	log.KV("tempflags", tempflags).Debug(colorized.Gray(loginfo))
	//先收集一次jackpot

	// 旋转空位 /////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 把需要转的槽位先集合起来,并且区分开，有效Bonus相邻的槽位放一个集合，其他槽位放一个集合
	neighborPos := []int{}
	unneighborPos := []int{}

	for index, val := range tempflags {
		if val != 0 {
			continue
		}
		fourPos := reels.AroundPos(index)
		// 判断上下左右位置 有bonus 把位置放入neighborPos,否则放入unneighbor
		n := false
		for _, p := range fourPos {
			if reels.PosCheck(p) {
				if tempflags[p] != 0 {
					neighborPos = append(neighborPos, index)
					n = true
					break
				}
			}
		}
		if !n {
			unneighborPos = append(unneighborPos, index)
		}
	}
	nebPos := len(neighborPos)
	randPos := append(neighborPos, unneighborPos...)
	log.KVs(log.Fields{"nebor": neighborPos, "unnebor": unneighborPos}).Debug("randpos")

	// 根据配置找出reel1-reel5分别需要几个
	respinConf := config_go.GetCherryRespin(int64(gameCount))
	if respinConf == nil {
		log.KV("gamecount", gameCount).ErrorStack(3, "not found respinConf")
		return "", -1, 0, 0
	}
	// 设置tempflags 里的空位reel1_5   reel1-reel5需要的数量+gameCount == 15
	for i, reel := range slottools.SlotReel1_5[int64(gameCount)] {
		num := reel[0]
		bneighbor := reel[1]
		for c := 0; c < int(num); c++ {
			randRange := 0
			// 配置需要保证优先处理完需要邻近找的reel
			if bneighbor == 1 {
				randRange = nebPos // 如果reel需要邻近格子，就优先从neighbor里找
			} else {
				randRange = len(randPos)
			}
			rd := rand.Intn(randRange)
			if rd < nebPos {
				nebPos -= 1
			}
			randindex := randPos[rd]
			randPos = append(randPos[:rd], randPos[rd+1:]...)
			tempflags[randindex] = int64(i) // 找到一个位置后，标记需要reel的编号 (i=0-4) -> reel1-reel5
		}
	}
	BonusMsg.BonusReelId = tempflags
	BonusMsg.Games = []*message.BonusGame{}
	// 循环respin /////////////////////////////////////////////////////////////////////////////////////////////////////
	// 根据tempflags的标记，进行respin
	optionReels, reelsnum := summaryOption(triggerType)
	BonusMsg.Option = reelsnum
	// 通过summary选择用option1 还是 option2
	for c := 0; c < gameCount; {
		msgGame := &message.BonusGame{BonusGridVal: []int64{}, BonusReels: []int64{}}
		// 单次respin逻辑
		bonusIncr := false
		over := true
		for i, flag := range tempflags {
			if flag > slottools.MAST_TYPE_SCALE {
				msgGame.BonusGridVal = append(msgGame.BonusGridVal, -1)
				continue
			}
			over = false
			if int(flag) >= len(optionReels) {
				log.KVs(log.Fields{"flag": flag, "len(optionReels)": len(optionReels), "reelOption": reelsnum}).
					ErrorStack(3, "flag >= len(optionReels)")
				return "", -1, 0, 0
			}
			ret, optionReelPos, err := slottools.Reel(optionReels[flag], 1, -1)
			if err != nil {
				log.KVs(log.Fields{"error": err, "flag": flag, "len(optionReels)": len(optionReels), "reelsnum": reelsnum, "optionReelPos": optionReelPos}).
					ErrorStack(3, "slottools.Reel error")
				return "", -1, 0, 0
			}
			msgGame.BonusGridVal = append(msgGame.BonusGridVal, optionReelPos)
			reels.Arr[i] = ret[0]
			if patternType(ret[0]) == slottools.BonusType {
				tempflags[i] = ret[0]
				for _, pos := range reels.AroundPos(i) {
					if reels.PosCheck(pos) {
						val := reels.Arr[pos]
						if patternType(val) == slottools.BonusType && val > slottools.BonusType*slottools.MAST_TYPE_SCALE {
							bonusIncr = true
						}
					}
				}
			}
		}
		msgGame.BonusReels = reels.Arr

		if over {
			break
		}

		log.Info(colorized.White(reels.String()))
		reset := false
		if bonusIncr {
			reset = true
			jp := CheckAvaildBonus(reels) // 处理需要+X1的bonus
			for k, v := range jp {
				jackpot[k] = jackpot[k] + v
			}
		}
		if b, jp := CheckExtra3Linked(reels); b { // 查找并处理如额外的三连bonus
			for k, v := range jp {
				jackpot[k] = jackpot[k] + v
			}
			reset = true
		}

		if reset {
			c = 0
			msgGame.ResetCount = true
		} else {
			c++
		}
		BonusMsg.Games = append(BonusMsg.Games, msgGame)
		log.KV("c", c).Debug(colorized.White(reels.String()))
	}

	full := true
	for _, v := range reels.Arr {
		if patternType(v) == slottools.BonusType && v > slottools.BaseBonus {
			odds := config_go.GetSlotCherryPattern(v).BOUNS_ODDS()
			BonusOdds += int64(odds)
		} else {
			full = false
		}
	}
	if full {
		BonusOdds = BonusOdds * 2
	}

	for k, count := range jackpot {
		jpConf := config_go.GetCherryJackPot(k)
		if jpConf == nil {
			log.KV("k", k).Error("not found GetCherryJackPot")
			return
		}
		if k <= slottools.JackpotMINOR {
			jackpotOdds += int64(jpConf.DftOdds()) * int64(count)
		} else {
			jackpotOdds += int64(tools.Randx_y(int(this.jackpot[k][0]), int(this.jackpot[k][1])))
		}
	}
	return combineId, rowId, BonusOdds, jackpotOdds
}

// 如果触发respin,处理需要+X1的bonus
func CheckAvaildBonus(reel *slottools.ReelsArray) (jp map[int64]int) {
	changemap := map[int]bool{} // 增加过的bonus标记一下，防止多次增加
	jp = map[int64]int{}
	var recur func(index int)
	recur = func(index int) {
		arrpos := reel.AroundPos(index)
		for _, pos := range arrpos {
			if reel.PosCheck(pos) {
				if reel.Arr[pos] == slottools.BaseBonus { // 图案是baseBonus
					reel.Arr[pos] = reel.Arr[pos] + 1
					changemap[pos] = true
					recur(pos)
				} else if patternType(reel.Arr[pos]) == slottools.BonusType && reel.Arr[pos] >= slottools.JackpotMINI { // 图案是jackpot
					jp[reel.Arr[pos]] = jp[reel.Arr[pos]] + 1
					reel.Arr[pos] = slottools.JackpotAlter
					changemap[pos] = true
					recur(pos)
				}
			}
		}
	}

	for i, v := range reel.Arr {
		if patternType(v) == slottools.BonusType && v > slottools.BonusType*slottools.MAST_TYPE_SCALE && !changemap[i] {
			reel.Arr[i] = reel.Arr[i] + 1
			if reel.Arr[i] > 310 { // 不能超过x10
				reel.Arr[i] = 310
			}
			recur(i)
		}
	}
	return jp
}

// 检查是否有额外3连的bonus
func CheckExtra3Linked(reel *slottools.ReelsArray) (bool, map[int64]int) {
	jp := map[int64]int{}
	var recur func(index int)
	recur = func(index int) {
		arrpos := reel.AroundPos(index)
		for _, pos := range arrpos {
			if reel.PosCheck(pos) {
				if reel.Arr[pos] == slottools.BaseBonus { // 图案是baseBonus,直接+1
					reel.Arr[pos] = slottools.BaseBonus + 1
					recur(pos)
				} else if patternType(reel.Arr[pos]) == slottools.BonusType && reel.Arr[pos] >= slottools.JackpotMINI { // 图案是jackpot, 收集后，变成X1
					jp[reel.Arr[pos]] = jp[reel.Arr[pos]] + 1
					reel.Arr[pos] = slottools.JackpotAlter
					recur(pos)
				}
			}
		}
	}

	ret := false
	for i, v := range reel.Arr {
		if v == slottools.BaseBonus {
			arrpos := reel.AroundPos(i)
			count := 1
			// 找3连，找到就break
			for _, pos := range arrpos {
				if reel.PosCheck(pos) && reel.Arr[pos] == slottools.BaseBonus {
					count++
					if count == 3 {
						break
					}
				}
			}

			// 如果满足3连，先把自己+1，然后递归的把周围相连的所有bonus+1
			if count >= 3 {
				ret = true
				reel.Arr[i] = slottools.BaseBonus + 1
				recur(i)
			}
		}
	}
	return ret, jp
}

// 通过权重拿option1 或 option2 轮盘
func summaryOption(t int) ([][]int64, int64) {
	rate1 := config_go.GetCherrySummary(int64(t)).RATE_1()
	rate2 := config_go.GetCherrySummary(int64(t)).RATE_2()
	rand := map[int64]int64{
		1: rate1,
		2: rate2,
	}
	v := common.RandomWeight32(rand)
	switch v {
	case 1:
		return slottools.SlotOption1, v
	case 2:
		return slottools.SlotOption2, v
	}
	log.KVs(log.Fields{"rand": rand, "v": v}).ErrorStack(3, "summaryOption random -1")
	return nil, 0
}

func patternType(patternId int64) int {
	return int(patternId) / slottools.MAST_TYPE_SCALE
}

// 回存slot数据
func (this *Slot) Save() {
	this.game.Cache().Save(this.RID, this.ModelSlot)
}

func (this *Slot) Model() i_mysql.IDbTable {
	return this.ModelSlot
}

func (this *Slot) Role() interfaces.IPlayer {
	return this.game.PlayerMgr().GetRoleByRID(this.RID)
}
