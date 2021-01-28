package interfaces

import "root/internal/common"

type IPackage interface {
	AddItem(itemId, count int64, reason common.REASON)
	SubItem(itemId, count int64, reason common.REASON)
}
