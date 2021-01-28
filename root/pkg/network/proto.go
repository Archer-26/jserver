package network

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"root/pkg/log"
	"strings"
)

type MsgTypeMap struct {
	typemap map[int32]protoreflect.MessageType // 此map初始化后，所有session都会并发读取
}

func NewMsgTypeMap() *MsgTypeMap {
	return &MsgTypeMap{
		typemap: make(map[int32]protoreflect.MessageType),
	}
}

func (this *MsgTypeMap) InitMsgParser(packageName, msgType string) *MsgTypeMap {
	enums, err := protoregistry.GlobalTypes.FindEnumByName(protoreflect.FullName(packageName + "." + msgType))
	if err != nil {
		log.KV("packageName", packageName).KV("msgType", msgType).Error("InitMsgParser error")
		panic(fmt.Errorf("InitMsgParser error packageName=%v msgType=%v", packageName, msgType))
		return this
	}

	values := enums.Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		msgTypeName := string(values.Get(i).Name())
		if strings.Contains(msgTypeName, "SEGMENT_BEGIN") || strings.Contains(msgTypeName, "SEGMENT_END") {
			continue
		}
		bodyName := handMsgType(msgTypeName)
		tp, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(packageName + "." + bodyName))
		if err == nil {
			this.typemap[int32(values.Get(i).Number())] = tp
		} else {
			log.KV("msgTypeName", msgTypeName).KV("bodyName", bodyName).Error("msg format error")
			panic(fmt.Errorf("msg format error msgTypeName=%v bodyName=%v", msgTypeName, bodyName))
		}
	}
	return this
}

func (this *MsgTypeMap) UnmarshalPbMsg(msgType int32, data []byte) proto.Message {
	tp, ok := this.typemap[msgType]
	if !ok {
		log.KV("msgId", msgType).KV("data", data).ErrorStack(5, "msg not regist parser")
		return nil
	}
	msg := tp.New().Interface().(proto.Message)
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.KVs(log.Fields{"msgType": msgType, "data": data}).Error("msg parse failed")
		return nil
	}
	return msg
}

func handMsgType(msgTypeName string) (name string) {
	words := strings.Split(msgTypeName, "_")
	for _, word := range words {
		if word == "C2S" || word == "S2C" ||
			word == "D2L" || word == "L2D" || // DB<->Login
			word == "G2L" || word == "L2G" || // Game<->Login
			word == "G2D" || word == "D2G" || // Game<->DB
			word == "GT2D" || word == "D2GT" || // Gate<->DB
			word == "GT2G" || word == "G2GT" || // Gate<->Game
			word == "GT2L" || word == "L2GT" || // Gate<->Login
			word == "L2ALL" || word == "GT2ALL" || // 广播消息
			word == "ALL2D" {
			name += word
		} else {
			name += strings.Title(strings.ToLower(word))
		}
	}
	return
}
