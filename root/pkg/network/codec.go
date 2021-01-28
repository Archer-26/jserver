package network

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type ICodec interface {
	Decode([]byte) ([][]byte, error)
	Encode(data []byte) []byte
	MaxDecode() int
}

type Codec struct {
}

var ErrRecvLen = errors.New("data is too long")

func (this *Codec) MaxDecode() int {
	return 1024 * 1024 * 5
}

type EchoCodec struct {
	Codec
	context bytes.Buffer
}

func (this *EchoCodec) Decode(data []byte) ([][]byte, error) {
	this.context.Write(data)
	var ret [][]byte = nil
	for {
		d, err := this.context.ReadBytes('\n')
		if err != nil {
			break
		}
		if len(ret) > this.MaxDecode() {
			return nil, ErrRecvLen
		}
		ret = append(ret, d[:len(d)-1])
	}
	return ret, nil
}

func (this *EchoCodec) Encode(data []byte) []byte {
	data = append(data, '\n')
	return data
}

type StreamCodec struct {
	Codec
	msglen  uint32
	context bytes.Buffer
}

const STREAM_HEADLEN = 4

func (this *StreamCodec) Decode(data []byte) ([][]byte, error) {
	this.context.Write(data)

	var ret [][]byte = nil
	for this.context.Len() >= STREAM_HEADLEN {
		if this.msglen == 0 {
			d := this.context.Bytes()
			this.msglen = binary.BigEndian.Uint32(d[:STREAM_HEADLEN])
			if int(this.msglen) > this.MaxDecode() {
				return nil, ErrRecvLen
			}
		}

		if int(this.msglen)+STREAM_HEADLEN > this.context.Len() {
			break
		}
		d := make([]byte, this.msglen+STREAM_HEADLEN)
		n, err := this.context.Read(d)
		if n != int(this.msglen)+STREAM_HEADLEN || err != nil {
			this.msglen = 0
			continue
		}
		this.msglen = 0
		ret = append(ret, d[STREAM_HEADLEN:])
	}
	return ret, nil
}

func (this *StreamCodec) Encode(data []byte) []byte {
	d := make([]byte, STREAM_HEADLEN)
	binary.BigEndian.PutUint32(d, uint32(len(data)))
	data = append(d, data...)
	return data
}
