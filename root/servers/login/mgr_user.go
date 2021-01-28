package login

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
	"math/rand"
	"root/internal/common"
	"root/internal/system"
	"root/pkg/abtime"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/models"
	"root/servers/login/logindb"
	"strings"
)

const MaxRoleSameServer = 1 //同一账号单服最大角色数
const SID = 1

var UserMgr = &UserManager{
	roles:     map[int64]*models.RoleInfo{},
	users:     map[int64]*models.ModelUser{},
	passports: map[string]*models.Passport{},
}

type UserManager struct {
	roles     map[int64]*models.RoleInfo  // map[RID]*models.RoleInfo
	users     map[int64]*models.ModelUser // map[UID]*models.ModelUser
	passports map[string]*models.Passport
	All_UID   []int64 // 可用的UID
}

func (this *UserManager) Login(passport string, ptype string, ip string, packageCode string) (role *models.RoleInfo, user *models.ModelUser, isNewUser bool) {
	key := PassportMerge(passport, ptype)
	pport := this.passports[key]
	if pport == nil || this.users[pport.UID] == nil { // 创建新用户
		user = this.newUser(passport, ptype, ip, SID)
		isNewUser = true
	} else {
		user = this.users[pport.UID]
	}

	if role = user.Roles[user.LastLoginRID]; role == nil {
		log.KVs(log.Fields{"LastLoginRID": user.LastLoginRID, "uid": user.UID, "len(player)": len(user.Roles)}).
			Error("not found player")
		if len(user.Roles) != 0 {
			for _, r := range user.Roles {
				role = r
				break
			}
		}
	}
	if role == nil {
		log.KV("UID", user.UID).Error("can not find player")
		return
	}

	user.LastLoginRID = role.RID
	user.LastLogin = abtime.Now().Unix()
	role.LastLogin = user.LastLogin
	user.LastSecurityCode = common.LoginSecurityCode(role.RID, user.UID, SID, isNewUser, user.LastLogin)
	user.PackageCode = packageCode
	return role, user, isNewUser
}

func (this *UserManager) newUser(passport, ptype string, ip string, sid int64) *models.ModelUser {
	if len(this.All_UID) == 0 {
		log.Error("UID is empty!!")
		return nil
	}
	uid := this.All_UID[len(this.All_UID)-1]
	this.All_UID = this.All_UID[:len(this.All_UID)-1]
	newUser := models.NewModelUser(uid, ip)
	this.users[newUser.UID] = newUser

	newPassport := &models.Passport{Type: ptype, PID: passport, UID: newUser.UID, CreateAt: newUser.CreatedAt, CreateIP: ip}
	key := PassportMerge(passport, ptype)
	this.passports[key] = newPassport
	newUser.Passports[key] = newPassport

	newRole := this.newRoleInfo(passport, ptype, ip, sid)
	this.roles[newRole.RID] = newRole
	newUser.LastLoginRID = newRole.RID
	newUser.Roles[newRole.RID] = newRole

	bytes, err := msgpack.Marshal(newUser)
	if err != nil {
		log.KV("error", err).Error("msgpack.Marshal(newUser) failed")
		return nil
	}

	// 新用户回存
	saveDB := network.NewPbMessage(&inner.L2DUserSave{Data: bytes}, inner.INNER_MSG_L2D_USER_SAVE.Int32())
	system.Send("", common.Login_Actor, logindb.LoginDBName, saveDB)

	log.KVs(log.Fields{
		"passport": passport,
		"type":     ptype,
		"ip":       ip,
		"rid":      newRole.RID,
		"uid":      newUser.UID,
	}).Info("newuser login")
	return newUser
}

func (this *UserManager) newRoleInfo(passport, ptype string, ip string, sid int64) *models.RoleInfo {
	pport := this.passports[PassportMerge(passport, ptype)]
	field := log.Fields{"passport": passport, "type": ptype, "ip": ip, "uid": pport.UID}
	user := this.users[pport.UID]
	if user == nil {
		log.KVs(field).Error("newrole user not find")
		return nil
	}

	if len(user.Roles) > MaxRoleSameServer {
		log.KVs(field).KV("MaxRoleSameServer", MaxRoleSameServer).
			Warn("newrole len(user.Roles) > MaxRoleSameServer")
		return nil
	}

	rid := int64(0)
	for i := 1; i < 9; i++ {
		calc_rid := common.UIDToRID(user.UID, i)
		if _, exist := user.Roles[calc_rid]; !exist {
			rid = calc_rid
			break
		}
	}
	if rid == 0 {
		log.KVs(log.Fields{"RID": user.Roles}).Error("newrole not find available for rid ")
		return nil
	}

	newRole := &models.RoleInfo{
		RID:      rid,
		UID:      user.UID,
		SID:      sid,
		CreateAt: abtime.Now().Unix(),
		CreateIP: ip,
	}
	newRole.Name = fmt.Sprintf("Player_%d", rid)
	user.Roles[newRole.RID] = newRole
	this.roles[newRole.RID] = newRole
	return newRole
}

func (this *UserManager) UserLen() int64 {
	return int64(len(this.roles))
}

func (this *UserManager) User(UID int64) *models.ModelUser {
	return this.users[UID]
}

func PassportMerge(passport, ptype string) string {
	return fmt.Sprintf("%s#_#%s", passport, ptype)
}

func PassportSplit(passport string) (string, string, bool) {
	strs := strings.Split(passport, "#_#")
	if len(strs) != 2 {
		return "", "", false
	}
	return strs[0], strs[1], true
}

// 开服后，必须在所有用户加载完成后调用，生成可用的所有UID，排除已有的
// globalNumb 服务器编号分服
func (this *UserManager) AllUID(globalNumb int32) {
	inival := int64(1010207)
	max := int64(8996219)

	used := 0
	all := make([]int64, 0, max-inival)
	for i := inival; i < max; i++ {
		val := int64(globalNumb)*10000000 + i
		if _, ok := this.users[val]; !ok {
			all = append(all, val)
			l := len(all)
			swappos := rand.Intn(l)
			all[swappos], all[l-1] = all[l-1], all[swappos]
		} else {
			used++
		}
	}
	log.KVs(log.Fields{"len": len(all), "be used count": used}).Info(colorized.Cyan("all uid build finish!"))
	this.All_UID = all
}
