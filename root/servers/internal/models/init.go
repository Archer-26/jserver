package models

import (
	"root/servers/internal/mysql/i_mysql"
)

func init() {
	i_mysql.RegisterTable(&ModelRole{})
	i_mysql.RegisterTable(&ModelItems{})
	i_mysql.RegisterTable(&ModelUser{})
	i_mysql.RegisterTable(&ModelSlot{})
	i_mysql.RegisterTable(&ModelBuildings{})
	i_mysql.RegisterTable(&ModelIsland{})
	i_mysql.RegisterTable(&ModelDice{})
}
