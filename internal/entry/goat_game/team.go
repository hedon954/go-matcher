package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

type Team struct {
	*glicko2.TeamBaseGlicko2
}

func CreateTeam(base *entry.TeamBase) (entry.Team, error) {
	t := &Team{}

	if err := t.withMatchStrategy(base); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Team) withMatchStrategy(base *entry.TeamBase) error {
	switch base.MatchStrategy {
	case constant.MatchStrategyGlicko2:
		t.TeamBaseGlicko2 = glicko2.CreateTeamBase(base)
	default:
		return fmt.Errorf("unsupported match strategy: %d", base.MatchStrategy)
	}
	return nil
}
