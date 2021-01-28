package login

import (
	"math/rand"
	"root/internal/common"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/log"
)

var GameMgr = &GameManager{
	gameActors: map[string]*gameInfo{},
}

type (
	GameManager struct {
		login      actor.IActor
		gameActors map[string]*gameInfo
	}

	// 游戏服基本信息
	gameInfo struct {
		Player_number int64 // 在线人数
	}
)

func (this *GameManager) Init(a actor.IActor) {
	this.login = a
	games := system.LocalActorFuzzy(common.IsGame)
	for _, game := range games {
		this.SetGameActor(game.GetID())
	}
}

func (this *GameManager) SetGameActor(actorId string) {
	if _, ok := this.gameActors[actorId]; !ok {
		log.KV("new game", actorId).Info("gameManager add game")
		this.gameActors[actorId] = &gameInfo{}
	}
}

// 返回game的actorId
func (this *GameManager) GetAvailableGame() string {
	// todo... 暂时随机找一个game
	arr := []string{}
	for actorid := range this.gameActors {
		arr = append(arr, actorid)
	}
	if len(arr) == 0 {
		return ""
	}
	return arr[rand.Intn(len(arr))]
}
