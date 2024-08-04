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
func (m *GroupMgr) CreateGroup(
	playerLimit int, mode constant.GameMode, modeVersion int64, strategy constant.MatchStrategy,
) (
	g entry.Group, err error,
) {
	base := entry.NewGroupBase(m.GenGroupID(), playerLimit, mode, modeVersion, strategy)
	switch mode {
	case constant.GameModeGoatGame:
		g, err = goat_game.CreateGroup(base)
	default:
		g = base
	}
	if err != nil {
		return nil, err
	}
	m.Add(g.ID(), g)
	return g, nil
}

func (m *GroupMgr) GenGroupID() int64 {
	return m.groupIDIter.Add(1)
}
