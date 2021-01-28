package network

import (
	"net"
	"reflect"
	"sync/atomic"
	"time"

	"root/pkg/log"
)

type OptionClient func(l *TcpClient)

func StartTcpClient(addr string, coder ICodec, handler INetHandler, op ...OptionClient) (INetClient, bool) {
	l := &TcpClient{
		addr:          addr,
		codec:         reflect.TypeOf(coder).Elem(),
		handler:       reflect.New(reflect.TypeOf(handler).Elem()).Interface().(INetHandler),
		maxread:       1024 * 10,
		extentionData: nil}

	for _, f := range op {
		f(l)
	}
	err := l.connect()
	if err == nil {
		return l, true
	}
	return l, false
}
func SetCSMaxRead(max uint64) OptionClient {
	return func(l *TcpClient) {
		l.maxread = max
	}
}
func SetCSExtra(extra interface{}) OptionClient {
	return func(l *TcpClient) {
		l.extentionData = extra
	}
}

type TcpClient struct {
	addr          string
	running       int32
	session       INetSession
	codec         reflect.Type
	handler       INetHandler
	reconnect     int
	maxread       uint64
	extentionData interface{}
}

func (this *TcpClient) connect() error {
	conn, err := net.DialTimeout("tcp", this.addr, time.Second)
	if err != nil {
		log.KVs(log.Fields{"remote": this.addr, "error": err}).Error("tcp connect failed")
		return err
	}
	session := newTcpSession(conn, reflect.New(this.codec).Interface().(ICodec), this.handler, this.maxread, this.extentionData)
	this.session = session
	return nil
}

func (this *TcpClient) Reconnect() bool {
	this.reconnect = 0
	for {
		last := this.session
		err := this.connect()
		if err == nil {
			if last != nil {
				last.CopyStore(this.session)
			}
			return true
		}
		this.reconnect++
		time.Sleep(time.Second * time.Duration(this.reconnect+3))
	}
}

func (this *TcpClient) IsRunning() bool         { return atomic.LoadInt32(&this.running) == 1 }
func (this *TcpClient) NetSession() INetSession { return this.session }
func (this *TcpClient) Stop() {
	if atomic.CompareAndSwapInt32(&this.running, 1, 0) {
		log.KV("addr", this.addr).Info("stop client")
		this.session.Stop()
	}
}
