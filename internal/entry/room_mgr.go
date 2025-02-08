package entry

import (
	"sync/atomic"

	"github.com/hedon954/go-matcher/pkg/collection"
)

type RoomMgr struct {
	*collection.Manager[int64, Room]
	roomIDIter atomic.Int64
}

// NewRoomMgr creates a room repository.
func NewRoomMgr(roomIDStart int64) *RoomMgr {
	mgr := &RoomMgr{
		Manager: collection.New[int64, Room](),
	}
	mgr.roomIDIter.Store(roomIDStart)
	return mgr
}
