package i_mysql

import "root/pkg/iniconfig"

var allDbTable = make(map[string]IDbTable)

func RegisterTable(dbTable IDbTable) {
	allDbTable[dbTable.ModelName()] = dbTable
}
func GetModelByName(tableName string) IDbTable {
	return allDbTable[tableName]
}

func RangeTableWithDbName(dbName string, f func(IDbTable) bool) {
	for _, dbTable := range allDbTable {
		for _, belongDb := range dbTable.BelongDB() {
			if dbName == iniconfig.String(belongDb) {
				if f(dbTable) {
					return
				}
			}
		}
	}
}

type IDbTable interface {
	BelongDB() []string
	ModelName() string
}
