package coder

import (
	"root/internal/common"
	"root/pkg/log"
	"root/pkg/network"
)

type MessageCoder struct {
	inner_typemap *network.MsgTypeMap
	typemap       *network.MsgTypeMap
}

func NewMessageCoder() *MessageCoder {
	return &MessageCoder{
		inner_typemap: network.NewMsgTypeMap().InitMsgParser("inner_message", "INNER_MSG"),
		typemap:       network.NewMsgTypeMap().InitMsgParser("message", "MSG"),
	}
}

func (s *MessageCoder) Encode(msg interface{}) []byte {
	ntMsg, ok := msg.(*network.Message)
	if !ok {
		log.ErrorStack(3, "data is not network.Message")
		return nil
	}
	return ntMsg.Buffer()
}
func (s *MessageCoder) Decode(bytes []byte) interface{} {
	msgId := network.Byte4ToUint32(bytes)
	var mm *network.MsgTypeMap
	if msgId > common.MESSAGE_SEGMENT {
		mm = s.typemap
	} else if msgId > common.INNER_MESSAGE_SEGMENT {
		mm = s.inner_typemap
	} else {
		return nil
	}
	return network.NewBytesMessageParse(bytes, mm)
}
