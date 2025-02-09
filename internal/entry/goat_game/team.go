package goat_game

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

type Team struct {
	*glicko2.TeamBaseGlicko2
}

func (t *Team) withMatchStrategy(base *entry.TeamBase, mgr *entry.GroupMgr) {
	t.TeamBaseGlicko2 = glicko2.CreateTeamBase(base, mgr)
}

func (t *Team) Encode() ([]byte, error) {
	return entry.Encode(t)
}

func (t *Team) Decode(data []byte) error {
	return entry.Decode(data, t)
}
