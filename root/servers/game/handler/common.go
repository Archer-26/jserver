package handler

import (
	"root/pkg/log"
	"root/servers/internal/message"
)

func LogResultHandler(rid int64, result message.MSG_RESULT) {
	if result != message.MSG_RESULT_SUCCESS {
		log.KVs(log.Fields{"result": result, "rid": rid}).InfoStack(2, "process failed result")
	}
}
