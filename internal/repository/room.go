package repository

import (
	"sync/atomic"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type RoomMgr struct {
	*collection.Manager[int64, entry.Room]
	roomIDIter atomic.Int64
}

// RoomMgr creates a room repository.
func NewRoomMgr(roomIDStart int64) *RoomMgr {
	mgr := &RoomMgr{
		Manager: collection.New[int64, entry.Room](),
	}
	mgr.roomIDIter.Store(roomIDStart)
	return mgr
}

func (m *RoomMgr) CreateRoom(t entry.Team) (r entry.Room, err error) {
	base := entry.NewRoomBase(m.roomIDIter.Add(1), t)
	switch t.Base().GameMode() {
	case constant.GameModeGoatGame:
		r, err = goat_game.CreateRoom(base, t)
	default:
		r = base
	}
	if err != nil {
		return nil, err
	}

	r.Base().AddTeam(t)
	m.Add(r.Base().ID(), r)
	return r, nil
}
