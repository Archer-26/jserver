package client

import (
	"math/rand"
	"root/pkg/actor"
	"root/pkg/log"
	"time"
)

type (
	Robots struct {
		conn_addrs []string
		actor.IActor

		robot []*Client
	}
)

func (this *Robots) Init(a actor.IActor) {
	this.IActor = a
	this.AddTimer(5*time.Second, 10, this.Update)
}

func (this *Robots) Stop() bool {
	return true
}
func (this *Robots) Update(dt int64) {
	// 每秒 随机n个玩家，在n秒内，登录，并进入游戏
	randCount := rand.Intn(1) + 1
	for randCount > 0 {
		randtime := rand.Intn(10) + 1
		randCount--
		log.KV("rand", randtime).Debug("addtime")
		this.AddTimer(time.Duration(randtime)*1000*time.Millisecond, -1, func(dt int64) {
			log.KV("rand", randtime).Debug("tick")
			//system.Regist(actor.NewActor(name, &Client{}, actor.SetMailBoxSize(100), actor.SetLocalized()))
		})
	}
}

func (this *Robots) ActorID() string {
	return this.GetID()
}

func (this *Robots) HandleMessage(actorMsg *actor.ActorMessage) {
	return
}
