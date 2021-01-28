package role

import (
	"root/internal/common"
	"root/internal/config/config_global"
	"root/internal/config/config_go"
	"root/pkg/abtime"
	"root/pkg/log"
	"root/pkg/network"
	"root/servers/game/event"
	"root/servers/game/interfaces"
	"root/servers/internal/message"
	"root/servers/internal/models"
	"root/servers/internal/mysql/i_mysql"
)

type (
	Role struct {
		RID  int64
		game interfaces.IGame
		*models.ModelRole
	}
)

func (this *Role) Init(rid int64, g interfaces.IGame) {
	this.RID = rid
	this.game = g
	this.ModelRole = models.NewModelRole(rid)
}

func (this *Role) BaseInfo(uid, sid int64, passport, passporttype, name string) {
	this.UID = uid
	this.SID = sid
	this.Passport = passport
	this.PassportType = passporttype
	this.Name = name
	this.CreatedAt = abtime.Seconds()
}

func (this *Role) Sync() {
	this.notifyModify()
}

func (this *Role) Save() {
	this.game.Cache().Save(this.RID, this.ModelRole)
}

func (this *Role) Model() i_mysql.IDbTable {
	return this.ModelRole
}

func (this *Role) Player() interfaces.IPlayer {
	return this.game.PlayerMgr().GetRoleByRID(this.RID)
}

func (this *Role) AddGold(count int64, itemLog *common.ItemLog) {
	this.Gold += count
	this.game.Cache().Save(this.RID, this.ModelRole)
	log.KVs(log.Fields{"rid": this.RID, "gold": this.Gold, "add": count, "itemLog": itemLog}).Info("AddGold")
	this.notifyModify()
}

func (this *Role) SubGold(count int64, itemLog *common.ItemLog) {
	this.Gold -= count
	if this.Gold < 0 {
		this.Gold = 0
	}
	this.game.Cache().Save(this.RID, this.ModelRole)
	log.KVs(log.Fields{"rid": this.RID, "gold": this.Gold, "count": count, "itemLog": itemLog}).Info("SubGold")
	this.notifyModify()
}

func (this *Role) AddExp(exp int64, itemLog *common.ItemLog) {
	this.Exp += exp
	this.game.Cache().Save(this.RID, this.ModelRole)
	log.KVs(log.Fields{"rid": this.RID, "addExp": exp, "curExp": this.Exp, "level": this.Level, "itemLog": itemLog}).Info("AddExp")
	this.notifyModify()

	var rewards = make(map[int32]int64)
	for {
		cfg := config_go.GetBasePlayerLevel(int64(this.Level + 1))
		if cfg == nil { //已经满级
			break
		}
		if this.Exp < int64(cfg.NEED_EXP()) { //不够升级
			break
		}

		config_global.GetRewardByPool(cfg.REWARD_POOL(), int64(this.Level), rewards)
		this.Exp = this.Exp - int64(cfg.NEED_EXP())
		this.Level += 1
		this.Player().AddItemMap(rewards, common.NewItemLogWithId(itemLog.Id, common.ROLE_LEVEL_UP))

		this.Player().Dispatch(&event.Ev_Role_Upgrade{Level: this.Level})
	}
}

func (this *Role) notifyModify() {
	slotpb := &message.NotifyRoleInfo{
		UID:      this.UID,
		RID:      this.RID,
		SID:      this.SID,
		Name:     this.Name,
		IconW:    "",
		Level:    this.Level,
		VipExp:   0,
		VipLevel: 0,
		Exp:      this.Exp,
		Gold:     this.Gold,
	}
	msg := network.NewPbMessage(slotpb, message.MSG_NOTIFY_ROLE_INFO.Int32())
	this.Player().SendMsg(msg)
}
