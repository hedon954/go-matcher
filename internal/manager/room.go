package manager

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type RoomMgr struct {
	*collection.Manager[int64, entry.Room]
}

// RoomMgr creates a room manager.
func NewRoomMgr() *RoomMgr {
	mgr := &RoomMgr{
		Manager: collection.New[int64, entry.Room](),
	}
	return mgr
}

func (m *TeamMgr) CreateRoom() (entry.Room, error) {
	// TODO
	return nil, nil
}
