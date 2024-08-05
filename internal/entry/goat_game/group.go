package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game/glicko2"
)

type Group struct {
	*entry.GroupBase
}

func CreateGroup(base *entry.GroupBase, p entry.Player) (entry.Group, error) {
	group := &Group{base}
	switch base.MatchStrategy {
	case constant.MatchStrategyGlicko2:
		return glicko2.CreateGroup(group, p.(*glicko2.Player)), nil
	default:
		return nil, fmt.Errorf("unsupported match strategy: %d", base.MatchStrategy)
	}
}
