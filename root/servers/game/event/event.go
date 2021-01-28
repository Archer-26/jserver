package event

import (
	"root/pkg/ev"
)

const (
	EV_ROLE_UPGRADE ev.EventType = iota // 升级事件
)

type (
	Ev_Role_Upgrade struct {
		Level int32
	}
)

func (s *Ev_Role_Upgrade) EType() ev.EventType {
	return EV_ROLE_UPGRADE
}
