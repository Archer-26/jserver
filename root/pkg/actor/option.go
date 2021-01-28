package actor

/*
	提供Actor修改默认值的接口，通过NewActor函数传参设置
*/

// 设置远端序列化实例
func Coder(coder ICoder) SystemOption {
	return func(system *ActorSystem) {
		system.remoteCoder = coder
	}
}

// 设置newlist最大值
func NewList(size int) SystemOption {
	return func(system *ActorSystem) {
		system.newList = make(chan IActor, 100)
	}
}

// 设置Actor监听的端口
func ActorAddr(addr string) SystemOption {
	return func(system *ActorSystem) {
		system.actorAddr = addr
	}
}

// 设置etcd监听的端口
func EtcdAddr(addr string) SystemOption {
	return func(system *ActorSystem) {
		system.etcdAddr = addr
	}
}

// 设置etcd 前缀
func EtcdPrefix(prefix string) SystemOption {
	return func(system *ActorSystem) {
		system.etcdPrefix = prefix
	}
}
