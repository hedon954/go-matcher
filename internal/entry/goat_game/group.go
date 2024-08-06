package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

type Group struct {
	*glicko2.GroupBaseGlicko2
}

func CreateGroup(base *entry.GroupBase) (entry.Group, error) {
	g := &Group{}

	if err := g.withMatchStrategy(base); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Group) withMatchStrategy(base *entry.GroupBase) error {
	switch base.MatchStrategy {
	case constant.MatchStrategyGlicko2:
		g.GroupBaseGlicko2 = glicko2.NewGroup(base)
	default:
		return fmt.Errorf("unsupported match strategy: %d", base.MatchStrategy)
	}
	return nil
}
