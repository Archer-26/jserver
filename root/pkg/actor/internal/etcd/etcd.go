package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"root/pkg/log"
	"root/pkg/log/colorized"
	"root/pkg/tools"
	"strings"
	"time"
)

const (
	ETCD_TIMEOUT   = 5 * time.Second //etcd连接超时
	ETCD_GRANT_TTL = 6               //etcd abtime to live
)

type (
	Etcd struct {
		endpoints string
		prefix    string

		clietcd *clientv3.Client
		lease   *clientv3.LeaseGrantResponse            // 租约
		alive   <-chan *clientv3.LeaseKeepAliveResponse // 续约及连接保活
		watch   clientv3.WatchChan                      // 监听etcd服务事件

		regist_in    chan func()
		newActor_out chan<- WatchEvent

		registedActor map[string]string // map[actorId]json 记录所有put 到etcd的kv
	}

	WatchEvent map[string]struct {
		T   mvccpb.Event_EventType
		Evt string
	}
)

func NewEtcd(endpoints, prefix string, event chan<- WatchEvent) *Etcd {
	return &Etcd{
		endpoints:     endpoints,
		prefix:        prefix,
		regist_in:     make(chan func(), 1000),
		registedActor: make(map[string]string),
		newActor_out:  event,
	}
}

// 初始化并启动etcd本地服务
// example: etcd := newEtcd(....).Startup()
func (this *Etcd) Startup() *Etcd {
	for !this.connect() {
		time.Sleep(time.Second)
	}
	this.run()
	return this
}

// 注册kv
func (this *Etcd) RegistEtcd(k, v string) {
	this.put2etcd(k, v)
}

func (this *Etcd) syncLocalInfo() {
	for k, v := range this.registedActor {
		this.put2etcd(k, v)
	}
}

//////////////////////////////////////////////// inner func ///////////////////////////////////////////
func (this *Etcd) connect() bool {
	var err error
	this.clietcd, err = clientv3.New(clientv3.Config{Endpoints: strings.Split(this.endpoints, "_"), DialTimeout: ETCD_TIMEOUT})
	if err != nil {
		log.KV("addr", this.endpoints).KV("err", err).Error("etcd create client failed")
		return false
	}

	this.lease, err = this.clietcd.Grant(context.TODO(), ETCD_GRANT_TTL)
	if err != nil {
		log.KV("addr", this.endpoints).KV("err", err).Error("etcd create lease failed")
		return false
	}

	this.alive, err = this.clietcd.KeepAlive(context.TODO(), this.lease.ID)
	if err != nil {
		log.KV("addr", this.endpoints).KV("err", err).Error("etcd keepalive failed")
		return false
	}

	// 断线重连的情况，连接成功后，重新把本地所有actor注册到etcd上
	this.syncLocalInfo()

	log.KV("addr", this.endpoints).KV("lease", this.lease.ID).Info(colorized.Magenta("etcd connect success!"))

	// 连接成功后，监听
	this.watch = this.clietcd.Watch(context.Background(), this.prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	return true
}

func (this *Etcd) put2etcd(k, v string) {
	if l, c := len(this.regist_in), cap(this.regist_in); l >= c*2/3 {
		log.KVs(log.Fields{"len": l, "cap": c, k: v}).Warn("registEtcd will full chan")
	}

	key := this.prefix + k
	this.regist_in <- func() {
		if _, err := this.clietcd.Put(context.TODO(), key, v, clientv3.WithLease(this.lease.ID)); err != nil {
			log.KV("addr", this.endpoints).KV("err", err).KV("key", key).KV("value", v).Error("etcd put failed")
		} else {
			log.KV("key", key).KV("value", v).Info("etcd put success")
			this.registedActor[k] = v
		}
	}
}

func (this *Etcd) run() {
	tools.GoEngine(func() {
		for {
			select {
			case resp := <-this.alive:
				if resp == nil {
					log.KV("addr", this.endpoints).Error("etcd keepalive failed resp == nil, start reconnect")
					for !this.connect() {
						log.KV("addr", this.endpoints).Warn(colorized.Magenta("etcd reconnect failed"))
						time.Sleep(time.Second)
					}
				}
			case fn := <-this.regist_in: // etcd put
				tools.Try(func() { fn() }, nil)
			case watchResp := <-this.watch:
				if watchResp.Err() != nil {
					log.KV("addr", this.endpoints).KV("err", watchResp.Err()).Error("etcd watch err")
					continue
				}
				evts := WatchEvent{}
				for _, e := range watchResp.Events {
					key, val := this.shiftStruct(e.Kv)
					evts[key] = struct {
						T   mvccpb.Event_EventType
						Evt string
					}{T: e.Type, Evt: val}
				}
				log.KV("evts", evts).Info(colorized.White("etcd watched a new event"))
				this.newActor_out <- evts // 交给remote处理
			}
		}
	})
}

func (this *Etcd) shiftStruct(kv *mvccpb.KeyValue) (k, v string) {
	k = string(kv.Key)
	v = string(kv.Value)
	k = strings.TrimLeft(k, this.prefix)
	return
}
