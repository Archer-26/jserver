package interfaces

type WholeRoleBytes map[string][]byte // [modelName]mspack
type IPlayerMgr interface {
	NewRole(gateSession string, RID, UID, SID int64, Passport, PassportType string) IPlayer
	RoleMod(gateSession string, RID int64, mod WholeRoleBytes) IPlayer
	RelateSession2Role(gateSession string, role IPlayer)

	GetRoleByRID(RID int64) IPlayer
	GetRoleBySession(gateSession string) IPlayer
	GetRoleByUID(UID int64) IPlayer
	GetRoleByAsync(RID int64, callback func(mod WholeRoleBytes))
	SetOfflineRole(gateSession string)

	CallbackFun(callbackId, RID int64, mod WholeRoleBytes)
}
