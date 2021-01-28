package interfaces

import "root/servers/internal/message"

type IDice interface {
	SaveGameInfo(msg *message.DiceSaveInfoReq)(res *message.DiceSaveInfoRes)
}

