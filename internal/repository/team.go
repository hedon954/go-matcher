package repository

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type TeamMgr struct {
	*collection.Manager[int64, entry.Team]
}

// NewTeamMgr creates a team repository.
func NewTeamMgr(groupIDStart int64) *TeamMgr {
	mgr := &TeamMgr{
		Manager: collection.New[int64, entry.Team](),
	}
	return mgr
}

func (m *TeamMgr) CreateTeam() (entry.Team, error) {
	// TODO
	return nil, nil
}
