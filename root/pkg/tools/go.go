package tools

import (
	"fmt"
	"os"
	"root/pkg/log"
	"runtime/debug"
	"time"
)

func Try(fn func(), catch func(ex interface{})) {
	defer func() {
		if r := recover(); r != nil {
			stack := fmt.Sprintf("[%v] panic recover %v", r, string(debug.Stack()))
			log.Error(stack)
			if catch != nil {
				catch(r)
			}
		}
	}()
	fn()
}

func GoEngine(fn func()) {
	go Try(func() { fn() }, nil)
}

//逻辑需求
func GoLogic(fn func()) {
	go Try(func() { fn() }, nil)
}

//特别说明 于$GOROOT/src/runtime/proc.go中加入此函数
//func GoroutineId() int64 {
//  _g_ := getg()
//  return _g_.goid
//}

//func GoroutineId() uint64 {
//	b := make([]byte, 64)
//	b = b[:runtime.Stack(b, false)]
//	b = bytes.TrimPrefix(b, []byte("goroutine "))
//	b = b[:bytes.IndexByte(b, ' ')]
//	n, _ := strconv.ParseUint(string(b), 10, 64)
//	return n
//}

func Exit(code int) {
	log.KV("code", code).FatalStack(1, "os.Exit")
	time.Sleep(time.Second)
	os.Exit(code)
}
