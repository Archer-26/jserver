package interfaces

import "root/servers/internal/mysql/i_mysql"

type (
	ICache interface {
		Save(rid int64, Imd i_mysql.IDbTable)
		Stop()
	}
)
