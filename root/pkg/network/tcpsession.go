package network

import (
	"bufio"
	"io"
	"net"
	"root/pkg/abtime"
	"root/pkg/log"
	"root/pkg/tools"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func newTcpSession(conn net.Conn, coder ICodec, handler INetHandler, maxread uint64, extentionData interface{}) *TcpSession {
	session := &TcpSession{
		id:      GenNetSessionId(),
		conn:    conn,
		running: 1,
		coder:   coder,
		handler: handler,
		recvQue: make(chan []byte, 16),
		sendQue: make(chan []byte, 16),
		maxread: maxread,
	}
	log.KVs(log.Fields{"id": session.Id(), "local": session.LocalAddr(), "remote": session.RemoteAddr()}).
		Info("new tcp session")
	handler.setSession(session)
	if extentionData != nil {
		session.StoreKV("ext", extentionData)
	}
	handler.OnSessionCreated()
	tools.GoEngine(func() { session.read() })
	tools.GoEngine(func() { session.write() })
	return session
}

type TcpSession struct {
	id      uint32
	conn    net.Conn
	storage sync.Map

	running int32
	coder   ICodec
	handler INetHandler
	recvQue chan []byte
	sendQue chan []byte
	maxread uint64
}

func (this *TcpSession) Type() SessionType    { return TYPE_TCP }
func (this *TcpSession) Id() uint32           { return this.id }
func (this *TcpSession) LocalAddr() net.Addr  { return this.conn.LocalAddr() }
func (this *TcpSession) RemoteAddr() net.Addr { return this.conn.RemoteAddr() }
func (this *TcpSession) RemoteIP() string {
	addr := this.RemoteAddr()
	switch v := addr.(type) {
	case *net.UDPAddr:
		if ip := v.IP.To4(); ip != nil {
			return ip.String()
		}
	case *net.TCPAddr:
		if ip := v.IP.To4(); ip != nil {
			return ip.String()
		}
	case *net.IPAddr:
		if ip := v.IP.To4(); ip != nil {
			return ip.String()
		}
	}
	return ""
}

func (this *TcpSession) IsRunning() bool                { return this != nil && atomic.LoadInt32(&this.running) == 1 }
func (this *TcpSession) StoreKV(key, value interface{}) { this.storage.Store(key, value) }
func (this *TcpSession) DeleteKV(key interface{})       { this.storage.Delete(key) }
func (this *TcpSession) Load(key interface{}) (value interface{}, ok bool) {
	return this.storage.Load(key)
}
func (this *TcpSession) CopyStore(other INetSession) {
	this.storage.Range(func(key, value interface{}) bool {
		other.StoreKV(key, value)
		return true
	})
}

func (this *TcpSession) Stop() {
	tools.Try(func() {
		if atomic.CompareAndSwapInt32(&this.running, 1, 0) {
			_ = this.conn.Close()
			close(this.recvQue)
			close(this.sendQue)
			log.KV("sessionId", this.id).Info("tcp session close")
		}
	}, nil)
}

func (this *TcpSession) SendMsg(msg []byte) {
	tools.Try(func() {
		if this.IsRunning() {
			this.sendQue <- msg
		}
	}, nil)
}

func (this *TcpSession) read() {
	reader := bufio.NewReader(this.conn)
	buffer := make([]byte, this.maxread)
	var err error
	var n int
	var datas [][]byte
	tools.Try(func() {
		for {
			if err = this.conn.SetReadDeadline(abtime.Now().Add(time.Second * 30)); err != nil {
				log.KVs(log.Fields{"sessionId": this.Id(), "err": err}).Info("tcp read SetReadDeadline")
				break
			}
			if n, err = reader.Read(buffer); err != nil {
				if err != nil || n == 0 {
					if operr, ok := err.(*net.OpError); ok && operr != nil {
						if operr.Err == syscall.EAGAIN || operr.Err == syscall.EWOULDBLOCK { //没数据了
							continue
						}
					}
					if err != io.EOF {
						if neterr, ok := err.(net.Error); ok && neterr.Timeout() { //调整"i/o timeout"日志级别
							log.KVs(log.Fields{"sessionId": this.Id(), "err": err}).Info("tcp read buff failed")
						} else {
							log.KVs(log.Fields{"sessionId": this.Id(), "err": err}).Warn("tcp read buff failed")
						}
					}
					break
				}
			}
			if datas, err = this.coder.Decode(buffer[:n]); err != nil {
				log.KVs(log.Fields{"sessionId": this.Id(), "err": err}).Error("tcp decode failed")
				break
			}
			for _, d := range datas {
				this.handler.OnRecv(d)
			}
		}
	}, nil)
	this.Stop()
	this.handler.OnSessionClosed()
}

func (this *TcpSession) write() {
	tools.Try(func() {
		for data := range this.sendQue {
			if err := this.conn.SetWriteDeadline(abtime.Now().Add(time.Second * 5)); err != nil {
				log.KVs(log.Fields{"sessionId": this.Id(), "err": err}).Info("tcp read SetWriteDeadline")
				break
			}
			_, err := this.conn.Write(this.coder.Encode(data))
			if err != nil {
				log.KVs(log.Fields{"sessionId": this.Id(), "err": err}).Info("tcp session write error")
				break
			}
		}
	}, nil)
	this.Stop()
}
