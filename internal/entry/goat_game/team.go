package goat_game

import (
	"encoding/gob"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

func init() {
	gob.Register(&Team{})
}

type Team struct {
	*glicko2.TeamBaseGlicko2
}

func CreateTeam(base *entry.TeamBase) entry.Team {
	t := &Team{}

	t.withMatchStrategy(base)
	return t
}

func (t *Team) withMatchStrategy(base *entry.TeamBase) {
	t.TeamBaseGlicko2 = glicko2.CreateTeamBase(base)
}

func (t *Team) Encode() ([]byte, error) {
	return entry.Encode(t)
}

func (t *Team) Decode(data []byte) error {
	return entry.Decode(data, t)
}
