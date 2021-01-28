package expect

import (
	"fmt"
	"root/pkg/log"
)

func Nil(v interface{},fields... log.Fields)  {
	if v != nil{
		msg := fmt.Sprintf("\n%v\n",v)
		for _,v := range fields{
			log.KVs(v).ErrorStack(3,msg)
		}
		panic(nil)
	}
}
func True(b bool,fields... log.Fields)  {
	if !b{
		msg := "assert false"
		for _,v := range fields{
			log.KVs(v).ErrorStack(3,msg)
		}
		panic(nil)
	}
}