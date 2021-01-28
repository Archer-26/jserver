package network

import (
	"fmt"
	"net"
	"reflect"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/tools"
	"strings"
	"sync/atomic"
)

type OptionListen func(l *TcpListener)

func StartTcpListen(addr string, coder ICodec, handler INetHandler, op ...OptionListen) INetListener {
	l := &TcpListener{
		addr:          addr,
		maxread:       1024 * 10,
		extentionData: nil}
	for _, f := range op {
		f(l)
	}
	if err := l.listen(coder, handler); err == nil {
		log.KV("addr", addr).Info(colorized.Blue("tcp listen success"))
		return l
	} else {
		log.KVs(log.Fields{"addr": addr, "error": err}).Error("tcp listen faild")
	}
	return nil
}
func SetSSMaxRead(max uint64) OptionListen {
	return func(l *TcpListener) {
		l.maxread = max
	}
}
func SetSSExtra(extra interface{}) OptionListen {
	return func(l *TcpListener) {
		l.extentionData = extra
	}
}

type TcpListener struct {
	addr          string
	listener      net.Listener
	running       int32
	ctype         reflect.Type
	htype         reflect.Type
	maxread       uint64
	extentionData interface{}
}

func (this *TcpListener) listen(coder ICodec, handler INetHandler) error {
	listener, err := net.Listen("tcp", this.addr)
	if err != nil {
		return err
	}
	this.listener = listener

	if atomic.CompareAndSwapInt32(&this.running, 0, 1) {
		this.ctype = reflect.TypeOf(coder).Elem()
		this.htype = reflect.TypeOf(handler).Elem()
		tools.GoEngine(func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					if !strings.Contains(err.Error(), "use of closed network connection") {
						log.KVs(log.Fields{"addr": this.addr, "err": err}).
							Warn("tcp accept error")
					}
					break
				}
				newHander := reflect.New(this.htype)
				newTcpSession(conn, reflect.New(this.ctype).Interface().(ICodec), newHander.Interface().(INetHandler), this.maxread, this.extentionData)
			}
			this.Stop()
		})
		return nil
	}
	return fmt.Errorf("listener is already running")
}

func (this *TcpListener) IsRunning() bool { return atomic.LoadInt32(&this.running) == 1 }
func (this *TcpListener) Stop() {
	if atomic.CompareAndSwapInt32(&this.running, 1, 0) {
		log.KV("addr", this.addr).Info("tcp stop listen")
		_ = this.listener.Close()
	}
}
