package mysql

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"root/pkg/actor"
	"root/pkg/iniconfig"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/servers/internal/mysql/i_mysql"
	"time"
)

/*
	Mysql 操作相关的acotr,不负责具体读写逻辑，通过闭包传入执行，关闭也由外部决定
*/
type (
	MySQLDatabase struct {
		actor.IActor
		DataBase     []*gorm.DB
		DataBaseName string
		Count        int
	}
)

func (this *MySQLDatabase) Init(a actor.IActor) {
	this.IActor = a
	this.DataBase = make([]*gorm.DB, this.Count, this.Count)
	addr := iniconfig.String("mysql_addr")

	log.KVs(log.Fields{
		"mysqlAddr":     addr,
		"databaseName":  this.DataBaseName,
		"databaseCount": this.Count,
	}).Info(colorized.Yellow("mysql Database init!"))
	for i := 0; i < this.Count; i++ {
		dn := fmt.Sprintf("%v%v", this.DataBaseName, i)
		if err := _createDatabase(dn, fmt.Sprintf(addr, "")); err != nil {
			panic(err)
		}

		this.DataBase[i] = _mysqlopen(fmt.Sprintf(addr, dn), 10, 10)
		if this.DataBase[i] == nil {
			log.Error("connect mysql failed")
			panic("connect mysql failed ")
		}
		this.DataBase[i].LogMode(true)
		i_mysql.RangeTableWithDbName(this.DataBaseName, func(iDbTable i_mysql.IDbTable) bool {
			if err := this.DataBase[i].Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(iDbTable).Error; err != nil {
				log.KVs(log.Fields{
					"model": iDbTable,
					"error": err,
				}).Error("Init_table, build table faild")
			} else {
				log.KVs(log.Fields{
					"table name": gorm.ToTableName(iDbTable.ModelName()),
					"model name": iDbTable.ModelName(),
				}).Info(colorized.Magenta("mysql model it's ok!"))
			}

			return false
		})

	}
}
func (this *MySQLDatabase) Stop() bool {
	return false // MySQLDatabase关闭，由外部调用
}

func (this *MySQLDatabase) ActorID() string {
	return this.GetID()
}

func (this *MySQLDatabase) HandleMessage(actorMsg *actor.ActorMessage) {
	return
}

func _createDatabase(dbname string, addr string) error {
	mysql, err := sql.Open("mysql", addr)
	defer mysql.Close()
	if err != nil {
		log.KV("database", dbname).KV("addr", addr).KV("err", err).Error("mysql connect failed")
		return err
	}

	buildDataBase_command := "CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"
	if _, err = mysql.Exec(fmt.Sprintf(buildDataBase_command, dbname)); err == nil {
		log.KV("database", dbname).KV("addr", addr).Info(colorized.Gray("mysql create success"))
	} else {
		log.KV("database", dbname).KV("addr", addr).KV("err", err).Error(colorized.Red("mysql.Exec"))
		return err
	}
	return nil
}

// 打开一个数据库
func _mysqlopen(addr string, maxopenconns, maxidleconns int) *gorm.DB {
	if addr == "" {
		log.Error("add is nil")
		return nil
	}

	db, err := gorm.Open("mysql", addr)
	if err != nil {
		log.KVs(log.Fields{"error": err, "addr": addr}).Error("open mysql error")
		return nil
	}

	db.DB().SetConnMaxLifetime(time.Second * 5)
	db.DB().SetMaxIdleConns(maxidleconns)
	db.DB().SetMaxOpenConns(maxopenconns)
	db.SingularTable(true) // gorm默认是复数形式，让grom转义struct名字的时候不用加上s
	return db
}
