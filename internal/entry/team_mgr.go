package entry

import (
	"sync/atomic"

	"github.com/hedon954/go-matcher/pkg/collection"
)

type TeamMgr struct {
	*collection.Manager[int64, Team]
	teamIDIter atomic.Int64
}

// NewTeamMgr creates a team repository.
func NewTeamMgr(teamIDStart int64) *TeamMgr {
	mgr := &TeamMgr{
		Manager: collection.New[int64, Team](),
	}
	mgr.teamIDIter.Store(teamIDStart)
	return mgr
}
