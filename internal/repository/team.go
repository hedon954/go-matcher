package repository

import (
	"fmt"
	"sync/atomic"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type TeamMgr struct {
	*collection.Manager[int64, entry.Team]
	teamIDIter atomic.Int64
}

// NewTeamMgr creates a team repository.
func NewTeamMgr(teamIDStart int64) *TeamMgr {
	mgr := &TeamMgr{
		Manager: collection.New[int64, entry.Team](),
	}
	mgr.teamIDIter.Store(teamIDStart)
	return mgr
}

func (m *TeamMgr) CreateTeam(g entry.Group) (t entry.Team, err error) {
	base := entry.NewTeamBase(m.teamIDIter.Add(1), g)

	switch base.GameMode {
	case constant.GameModeGoatGame:
		t = goat_game.CreateTeam(base)
	case constant.GameModeTest:
		t = base
	default:
		return nil, fmt.Errorf("unsupported game mode: %d", base.GameMode)
	}

	t.Base().AddGroup(g)

	// NOTE: don't add team to manager here.
	// because it may be created in match process for temp,
	// only add it after match success.
	// m.Add(t.Base().ID(), t)
	return t, nil
}

func (m *TeamMgr) CreateAITeam(g entry.Group) (t entry.Team, err error) {
	base := entry.NewTeamBase(m.teamIDIter.Add(1), g)
	base.Base().RemoveGroup(g.ID())
	return t, nil
}
