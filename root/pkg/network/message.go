package network

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"root/pkg/log"
)

type Message struct {
	msgId int32
	pb    proto.Message
	bytes []byte
}

func (this *Message) Buffer() []byte {
	if len(this.bytes) > 0 {
		return this.bytes
	}
	if this.pb != nil {
		data, err := proto.Marshal(this.pb)
		if err != nil {
			log.KV("err", err).ErrorStack(3, "marshal pb failed")
			return nil
		}
		return append(Uint32ToByte4(uint32(this.msgId)), data...)
	} else {
		return Uint32ToByte4(uint32(this.msgId))
	}
}

func (this *Message) MsgId() int32         { return this.msgId }
func (this *Message) Proto() proto.Message { return this.pb }

func (this *Message) parse(data []byte, mm *MsgTypeMap) *Message {
	dlen := len(data)
	if dlen < 4 {
		log.KV("data", data).Error("err msg length")
		return nil
	}

	this.msgId = int32(Byte4ToUint32(data[:4]))
	pb := mm.UnmarshalPbMsg(this.msgId, data[4:])
	if nil == pb {
		return nil
	} else {
		this.pb = pb
		this.bytes = data
	}
	return this
}

// 业务逻辑层主要使用接口
func NewPbMessage(pb proto.Message, msgId int32) *Message {
	msg := &Message{msgId: msgId, pb: pb}
	return msg
}

// gate中转消息主要使用接口，避免gate做无意义的序列化工作
func NewBytesMessage(data []byte) *Message {
	msg := &Message{msgId: int32(Byte4ToUint32(data[:4])), bytes: data, pb: nil}
	return msg
}

// 远端actor通信主要使用接口
func NewBytesMessageParse(data []byte, mm *MsgTypeMap) *Message {
	msg := &Message{}
	return msg.parse(data, mm)
}

func Byte4ToUint32(data []byte) (result uint32) {
	buff := bytes.NewBuffer(data)
	binary.Read(buff, binary.BigEndian, &result)
	return
}

func Uint32ToByte4(data uint32) (result []byte) {
	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, data)
	return buff.Bytes()
}
