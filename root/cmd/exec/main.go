package main

import (
	"flag"
	"math/rand"
	"root/internal/common"
	"root/internal/config/config_go"
	"root/internal/system"
	"root/pkg/actor"
	"root/pkg/expect"
	"root/pkg/iniconfig"
	"root/pkg/log"
	"root/pkg/tools"
	"root/servers/client"
	"root/servers/game"
	"root/servers/gate"
	"root/servers/login"
	"time"
)

func main() {
	tools.Try(func() {
		initConf()
		err := config_go.Load(iniconfig.String("configjson"))
		expect.Nil(err,log.Fields{"info":"load config error!"})

		common.RefactorConfig()
		system.InitActorSystem()
		switch iniconfig.AppType() {
		case common.Gate_Actor:
			system.Regist(actor.NewActor(common.GateName(iniconfig.AppId()), &gate.Gate{}))
			// 如果进程是网关，特殊处理remoteCoder 只转发消息，不序列化
			system.GateSpecificRemoteCoder()
		case common.Login_Actor:
			system.Regist(actor.NewActor(common.Login_Actor, &login.Login{}))
		case common.Game_Actor:
			system.Regist(actor.NewActor(common.GameName(iniconfig.AppId()), &game.Game{}))

		case "Client":
			system.Regist(actor.NewActor("Client", &client.Client{}, actor.SetLocalized()))
		case "Robot":
			system.Regist(actor.NewActor("Robot", &client.Robots{}, actor.SetLocalized(), actor.AvailWheel()))
		case "All":
			system.Regist(actor.NewActor(common.GateName(1), &gate.Gate{Parse: true}))
			system.Regist(actor.NewActor(common.Login_Actor, &login.Login{}))
			system.Regist(actor.NewActor(common.GameName(1), &game.Game{}))
		}
		system.Startup()
	}, nil)
	time.Sleep(time.Second * 2) // 等一会，让日志写完
}

func initConf() {
	rand.Seed(time.Now().UTC().Unix()) //设置随机数种子
	tools.Try(func() {
		registerFlag()                    //解析命令参数
		appName, appId := loadINIConfig() //服务器配置初始化
		common.InitLog(int32(flag.Lookup("log").Value.(flag.Getter).Get().(int)), appName, appId)
		if pprofPort := iniconfig.Int32("pprof_port"); pprofPort > 0 {
			tools.PProfInit(pprofPort)
		}
	}, func(ex interface{}) {
		time.Sleep(time.Second * 1)
	})
}

func registerFlag() {
	flag.String("ini", "", "ini file path")
	flag.String("app", "", "app type")
	flag.String("config", "", "config table path")
	flag.Int("id", 0, "app id")
	flag.Int("log", 0, "log level, if debug log=0")
	flag.Int("gm", 0, "open gm=1, default close gm=0")

	flag.Parse()
}

func loadINIConfig() (appType string, appid int) {
	appName := flag.Lookup("app")
	appId := flag.Lookup("id")
	configPath := flag.Lookup("ini")
	err := iniconfig.LoadINIConfig(appName.Value.String(), appId.Value.(flag.Getter).Get().(int), configPath.Value.String())
	expect.Nil(err,log.Fields{"configPath": configPath.Value.String()})
	return appName.Value.String(), appId.Value.(flag.Getter).Get().(int)
}
