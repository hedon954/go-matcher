package repository

import (
	"fmt"
	"sync/atomic"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game"
	"github.com/hedon954/go-matcher/internal/entry/test_game"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type RoomMgr struct {
	*collection.Manager[int64, entry.Room]
	roomIDIter atomic.Int64
}

// NewRoomMgr creates a room repository.
func NewRoomMgr(roomIDStart int64) *RoomMgr {
	mgr := &RoomMgr{
		Manager: collection.New[int64, entry.Room](),
	}
	mgr.roomIDIter.Store(roomIDStart)
	return mgr
}

func (m *RoomMgr) CreateRoom(teamLimit int, t entry.Team) (r entry.Room, err error) {
	base := entry.NewRoomBase(m.roomIDIter.Add(1), teamLimit, t)

	switch base.GameMode {
	case constant.GameModeTest:
		r = test_game.CreateRoom(base)
	case constant.GameModeGoatGame:
		r = goat_game.CreateRoom(base)
	default:
		return nil, fmt.Errorf("unsupported game mode: %d", base.GameMode)
	}

	r.Base().AddTeam(t)

	// NOTE: don't add room to manager here.
	// because it may be created in match process for temp,
	// only add it after match success.
	// m.Add(r.Base().ID(), r)
	return r, nil
}
