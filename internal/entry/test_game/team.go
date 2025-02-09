package test_game

import (
	"encoding/gob"

	"github.com/hedon954/go-matcher/internal/entry"
)

func init() {
	gob.Register(&Team{})
}

type Team struct {
	*entry.TeamBase
}

func CreateTeam(base *entry.TeamBase) entry.Team {
	return &Team{
		TeamBase: base,
	}
}

func (t *Team) Encode() ([]byte, error) {
	return entry.Encode(t)
}

func (t *Team) Decode(data []byte) error {
	return entry.Decode(data, t)
}
