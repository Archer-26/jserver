package coder

import (
	"root/internal/common"
	"root/pkg/log"
	"root/pkg/network"
)

// 专门给网关用，网关不需要(反)序列化用户消息
type BytesCoder struct {
	inner_typemap *network.MsgTypeMap
}

func NewBytesCoder() *BytesCoder {
	return &BytesCoder{
		inner_typemap: network.NewMsgTypeMap().InitMsgParser("inner_message", "INNER_MSG"),
	}
}

func (s *BytesCoder) Encode(msg interface{}) []byte {
	ntMsg, ok := msg.(*network.Message)
	if !ok {
		log.ErrorStack(3, "data is not network.Message")
		return nil
	}
	return ntMsg.Buffer()
}
func (s *BytesCoder) Decode(bytes []byte) interface{} {
	msgId := network.Byte4ToUint32(bytes)
	if msgId > common.MESSAGE_SEGMENT {
		return network.NewBytesMessage(bytes)
	} else if msgId > common.INNER_MESSAGE_SEGMENT {
		return network.NewBytesMessageParse(bytes, s.inner_typemap)
	}
	log.KVs(log.Fields{
		"msgId":                        msgId,
		"common.MESSAGE_SEGMENT":       common.MESSAGE_SEGMENT,
		"common.INNER_MESSAGE_SEGMENT": common.INNER_MESSAGE_SEGMENT,
	}).ErrorStack(4, "msgid not in of range ")
	return nil
}
