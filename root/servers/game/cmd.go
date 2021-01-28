package game

import (
	"root/internal/common"
	"root/internal/system"
	"root/pkg/log"
	"strconv"
)

func (this *Game) initCmd() {
	system.RegistCmd(this.GetID(), "addgold", this.CMD_AddGold)
}

//p1 rid
//p2 count
func (this *Game) CMD_AddGold(param []string) {
	if len(param) != 2 {
		return
	}

	rid, _ := strconv.Atoi(param[0])
	count, _ := strconv.Atoi(param[1])
	role := this.playerMgr.GetRoleByRID(int64(rid))
	role.AddItemMap(map[int32]int64{common.GoldID: int64(count)}, common.NewItemLog(common.CMD))
	log.KVs(log.Fields{"RID": rid, "gold": role.Gold()}).Info("add gold success!")
}
