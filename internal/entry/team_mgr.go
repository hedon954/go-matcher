package entry

import (
	"sync/atomic"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/log"
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

// Encode encodes all teams into a map of game modes to their encoded bytes.
//
//nolint:dupl
func (m *TeamMgr) Encode() map[constant.GameMode][][]byte {
	res := make(map[constant.GameMode][][]byte, m.Len())
	m.Range(func(id int64, t Team) bool {
		bs, err := t.Encode()
		if err != nil {
			log.Error().Any("team", t).Err(err).Msg("failed to encode team")
			return true
		}
		res[t.Base().GameMode] = append(res[t.Base().GameMode], bs)
		return true
	})
	return res
}
