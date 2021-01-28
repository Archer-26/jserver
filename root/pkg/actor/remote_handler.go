package actor

import (
	"github.com/golang/protobuf/proto"
	"root/pkg/abtime"
	"root/pkg/actor/internal/actorpb/protofile"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/network"
)

// 心跳ping
func heatbeatpingHander(s *actorSessionHandler, msg proto.Message) {
	recv := msg.(*protofile.ActorHeatbeatPing)
	netMsg := network.NewPbMessage(&protofile.ActorHeatbeatPong{Addr: recv.Addr}, int32(protofile.MSG_TYPE_ACTOR_HEATBEAT_PONG))
	s.INetSession.SendMsg(netMsg.Buffer())
}

// 心跳pong
func heatbeatpongHander(s *actorSessionHandler, msg proto.Message) {
	recv := msg.(*protofile.ActorHeatbeatPong)
	cli, ok := s.addr2session[recv.Addr]
	if !ok {
		log.KV("addr", recv.Addr).Error("heatbeatpongHander,not found addr")
		return
	}
	cli.heatbeatTime = abtime.Milliseconds()
}

// 注册远端actor
func registHander(s *actorSessionHandler, msg proto.Message) {
	recv := msg.(*protofile.ActorRegist)
	for _, actorId := range recv.ActorId {
		s.setActor2Session(actorId, s.Id()) // 远端连接本端成功，发送信息 映射actor->session
		log.KV("source", actorId).Info(colorized.Cyan("remote registered request"))
	}
}

// 接收远端acotr 消息
func recvHander(s *actorSessionHandler, msg proto.Message) {
	recv := msg.(*protofile.ActorMessage)
	d := s.Coder.Decode(recv.Data)
	var message interface{}
	if d != nil {
		message = d.(*network.Message)
	}

	s.remoteMsg <- &ActorMessage{
		sourceId:    recv.SourceId,
		targetId:    recv.TargetId,
		gateSession: recv.GateSession,
		msgId:       recv.MsgId,
		data:        message,
	}
}
