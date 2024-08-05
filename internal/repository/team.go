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
	switch g.Base().GameMode {
	case constant.GameModeGoatGame:
		t, err = goat_game.CreateTeam(base, g)
	case constant.GameModeTest:
		t = base
	default:
		return nil, fmt.Errorf("unsupported game mode: %d", g.Base().GameMode)
	}
	if err != nil {
		return nil, err
	}

	t.Base().AddGroup(g)
	m.Add(t.Base().ID(), t)
	return t, nil
}
