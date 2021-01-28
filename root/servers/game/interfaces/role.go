package interfaces

type IRole interface {
	BaseInfo(uid, sid int64, passport, passporttype, name string)
}
