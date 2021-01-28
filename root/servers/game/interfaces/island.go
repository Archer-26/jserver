package interfaces

import "root/servers/internal/message"

type IIsland interface {
	Place(island *message.Island)
	CheckpointFinish(island *message.CheckpointFinishReq)
	SaveJsonInfo(island *message.SaveIslandInfoReq)
}
