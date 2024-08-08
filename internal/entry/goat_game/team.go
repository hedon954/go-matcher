package goat_game

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

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
