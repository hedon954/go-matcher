package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game/glicko2"
)

type Team struct {
	*entry.TeamBase
}

func CreateTeam(base *entry.TeamBase, g entry.Group) (entry.Team, error) {
	team := &Team{base}
	switch g.Base().MatchStrategy {
	case constant.MatchStrategyGlicko2:
		return glicko2.CreateTeam(team, g.(*glicko2.Group)), nil
	default:
		return nil, fmt.Errorf("unsupported match strategy: %d", g.Base().MatchStrategy)
	}
}
