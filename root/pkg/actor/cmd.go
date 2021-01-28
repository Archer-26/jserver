package actor

import (
	"bufio"
	"fmt"
	"os"
	"root/pkg/actor/internal/actorpb/protofile"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
	"root/pkg/tools"
	"strings"
	"sync"
)

type (
	instructs struct {
		sys     *ActorSystem
		actorId string
		f       func([]string)
	}
)

var commands = make(map[string]instructs)
var lock sync.RWMutex

func InitCmd() {
	commands["h"] = instructs{actorId: "", f: func(i []string) {
		info := fmt.Sprintf("\ninstructions:\n")
		for cmd, _ := range commands {
			info += fmt.Sprintf("%v \n", cmd)
		}
		info += "\n"
		log.Info(info)
	},
		sys: nil}

	go loop()
}

func loop() {
	for {
		tools.Try(func() {
			reader := bufio.NewReader(os.Stdin)
			result, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("read error:", err)
			}
			result = result[:len(result)-1]
			args := strings.Split(result, " ")
			cmdName := args[0]
			lock.RLock()
			ins, exist := commands[cmdName]
			lock.RUnlock()
			if !exist {
				log.KV("cmdName", cmdName).Info("无效的命令\r\n")
			} else {
				param := args[1:]
				if ins.actorId == "" || ins.sys == nil {
					ins.f(param)
				} else {
					if ins.sys != nil {
						ins.sys.LocalSend("", ins.actorId, "CMD", func() {
							ins.f(param)
						})
					} else {
						log.KV("cmd", cmdName).Error("ins.sys == nil")
					}
				}
			}
		}, nil)
	}
}

func (s *ActorSystem) RegistCmd(actorId, cmd string, f func([]string)) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := commands[cmd]; ok {
		log.KV("cmd", cmd).ErrorStack(3, "repeated cmd")
		return
	}
	commands[cmd] = instructs{actorId: actorId, f: f, sys: s}
}

func (s *ActorSystem) actorInfo(param []string) {
	log.Info(colorized.Yellow("--------------------------------- local actor -------------------------------------"))
	s.actorCache.Range(func(key, value interface{}) bool {
		info := fmt.Sprintf("actor:%v mail-Box len:%v", key, len(value.(*actor).mailBox))
		log.Info(colorized.Red(info))
		return true
	})

	s.remoteMgr.sfun <- func() { // cmd指令
		log.Info(colorized.Yellow("--------------------------------- remote actor --------------------------------"))
		for actorId, sessionid := range s.remoteMgr.actor2session {
			sess := s.remoteMgr.sessions[sessionid]
			addr := "<---client"
			if sess != nil {
				addr = sess.Addr
			} else {
				addr = "not found"
			}
			info := fmt.Sprintf("actor:%v sessionId:%v addr:%v", actorId, sessionid, addr)
			log.Info(colorized.Red(info))
		}
		log.Info("")
		log.Info(colorized.Yellow("-------------------------------------------------------------------------------"))
	}
}

func (s *ActorSystem) detect(param []string) {
	if len(param) != 2 {
		return
	}
	source := param[0]
	target := param[1]

	if _, ok := s.actorCache.Load(source); !ok {
		log.KV("source", source).Red().Warn("source:%v is not local actor")
		return
	}
	s.Send("", source, target, network.NewPbMessage(nil, int32(protofile.MSG_TYPE_ACTOR_DETECT)))
}
