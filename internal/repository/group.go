package repository

import (
	"sync/atomic"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type GroupMgr struct {
	*collection.Manager[int64, entry.Group]

	// groupIDIter is a counter used to generate unique group IDs.
	// Before we shut down the server, we need to store the current groupID,
	// and then restart the server from the stored groupID.
	// This is done to avoid conflicts with group IDs that were generated before the server was shut down.
	groupIDIter atomic.Int64
}

// NewGroupMgr creates a group repository, `groupIDStart`: the starting group ID.
func NewGroupMgr(groupIDStart int64) *GroupMgr {
	mgr := &GroupMgr{
		Manager: collection.New[int64, entry.Group](),
	}
	mgr.groupIDIter.Store(groupIDStart)
	return mgr
}

// CreateGroup creates a group according to `pto.PlayerInfo`.
func (m *GroupMgr) CreateGroup(playerLimit int, p entry.Player) (
	g entry.Group, err error,
) {
	base := entry.NewGroupBase(m.GenGroupID(), playerLimit, p.Base())
	switch p.Base().GameMode {
	case constant.GameModeGoatGame:
		g, err = goat_game.CreateGroup(base, p)
	default:
		g = base
	}
	if err != nil {
		return nil, err
	}
	_ = g.Base().AddPlayer(p)
	m.Add(g.ID(), g)
	return g, nil
}

func (m *GroupMgr) GenGroupID() int64 {
	return m.groupIDIter.Add(1)
}
