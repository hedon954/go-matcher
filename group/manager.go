package group

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/hedon954/go-matcher/common"
	"github.com/hedon954/go-matcher/enum"
	"github.com/hedon954/go-matcher/internal/gamemode/goatgame"
	"github.com/hedon954/go-matcher/player"
	"github.com/hedon954/go-matcher/pto"
)

type Manager struct {
	playerMgr *player.Manager

	sync.RWMutex
	groups map[int64]common.Group

	// groupIDIter is a counter used to generate unique group IDs.
	// Before we shut down the server, we need to store the current groupID,
	// and then restart the server from the stored groupID.
	// This is done to avoid conflicts with group IDs that were generated before the server was shut down.
	groupIDIter atomic.Int64
}

// NewManager creates a group manager, `groupIDStart`: the starting group ID.
func NewManager(groupIDStart int64) *Manager {
	mgr := &Manager{
		groups: make(map[int64]common.Group, 1024),
	}
	mgr.groupIDIter.Store(groupIDStart)
	return mgr
}

// CreateGroup creates a group according to `pto.PlayerInfo`.
func (m *Manager) CreateGroup(info *pto.PlayerInfo) (common.Group, error) {
	groupID := m.GenGroupID()

	playerBase := &common.PlayerBase{
		Uid:               info.Uid,
		GroupID:           groupID,
		GameMode:          info.GameMode,
		ModeVersion:       info.ModeVersion,
		MatchStrategy:     info.MatchStrategy,
		UnityNamespace:    "TODO",
		UnityNamespacePre: info.UnityNamespacePre,
	}
	playerBase.SetVoiceState(common.PlayerVoiceStateOff)
	playerBase.SetOnlineState(common.PlayerOnlineStateGroup)

	var p common.Player
	var g common.Group
	var err error

	switch info.GameMode {
	case enum.GameModeGoat:
		if p, err = goatgame.CreatePlayer(playerBase); err != nil {
			break
		}
		if g, err = goatgame.CreateGroup(groupID, p); err != nil {
			break
		}
	default:
		return nil, fmt.Errorf("unsupported game mode: %d", info.MatchStrategy)
	}

	if err != nil {
		return nil, err
	}

	g.Base().SetState(common.GroupStateInvite)

	m.playerMgr.AddPlayer(p)
	m.AddGroup(groupID, g)
	return g, nil
}

func (m *Manager) GetGroup(id int64) common.Group {
	m.RLock()
	defer m.RUnlock()
	return m.groups[id]
}

func (m *Manager) GroupExists(id int64) bool {
	m.RLock()
	defer m.RUnlock()
	return m.groups[id] == nil
}

func (m *Manager) AddGroup(id int64, group common.Group) {
	m.Lock()
	defer m.Unlock()
	m.groups[id] = group
}

func (m *Manager) DelGroup(id int64) {
	m.Lock()
	defer m.Unlock()
	delete(m.groups, id)
}

func (m *Manager) GenGroupID() int64 {
	return m.groupIDIter.Add(1)
}
