package cachemod

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
	"root/internal/system"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/servers/game/interfaces"
	"root/servers/internal/inner_message/inner"
	"root/servers/internal/mysql/i_mysql"
	"time"
)

/*
	每个Actor 需要一个CacheModels 各自管理各自线程的缓存Model数据
*/

type (
	CacheModels struct {
		game  interfaces.IGame
		cache map[int64]map[string]i_mysql.IDbTable // [rid][modelName]*models.xxxx
	}
)

func Init(a interfaces.IGame) *CacheModels {
	cache := &CacheModels{
		game:  a,
		cache: make(map[int64]map[string]i_mysql.IDbTable),
	}

	a.AddTimer(5*time.Second, -1, cache.update)
	system.RegistCmd(a.GetID(), fmt.Sprintf("%v_cache", a.GetID()), func([]string) {
		fmt.Printf("cache size:%v\n", len(cache.cache))
	})
	return cache
}

func (this *CacheModels) Stop() {
	log.KV("actor", this.game.GetID()).
		KV("the amount of player that need to be saved", len(this.cache)).
		Info(colorized.Magenta("CacheModels Stop"))
	this.update(0)

	// 发送一条结束消息，表示game不在回存数据
	msg := network.NewPbMessage(&inner.G2DGameStop{}, inner.INNER_MSG_G2D_GAME_STOP.Int32())
	system.Send("", this.game.GetID(), this.game.DBName(), msg)
}

func (this *CacheModels) Save(rid int64, Imd i_mysql.IDbTable) {
	mdname := Imd.ModelName()
	m, exist := this.cache[rid]
	if !exist {
		this.cache[rid] = make(map[string]i_mysql.IDbTable)
		m = this.cache[rid]
	}
	m[mdname] = Imd
}

func (this *CacheModels) update(dt int64) {
	for rid, models := range this.cache {
		for name, imd := range models {
			bytes, err := msgpack.Marshal(imd)
			if err != nil {
				log.KVs(log.Fields{"error": err.Error(), "modName": name}).ErrorStack(3, " msgpack.Marshal(mod)")
				return
			}
			msg := network.NewPbMessage(&inner.G2DModelSave{Data: bytes, ModelName: imd.ModelName(), RID: rid}, inner.INNER_MSG_G2D_MODEL_SAVE.Int32())
			system.Send("", this.game.GetID(), this.game.DBName(), msg)
			log.KVs(log.Fields{"name": name, "rid": rid}).Info("cache models update")
		}
		delete(this.cache, rid)
	}
}
