package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Player struct {
	*glicko2.PlayerBaseGlicko2
}

func CreatePlayer(base *entry.PlayerBase, pInfo *pto.PlayerInfo) (entry.Player, error) {
	p := &Player{}

	if err := p.withMatchStrategy(base, pInfo.Glicko2Info); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Player) withMatchStrategy(base *entry.PlayerBase, info *pto.Glicko2Info) error {
	switch base.MatchStrategy {
	case constant.MatchStrategyGlicko2:
		p.PlayerBaseGlicko2 = glicko2.CreatePlayerBase(base, info)
	default:
		return fmt.Errorf("unknown match strategy: %d", base.MatchStrategy)
	}
	return nil
}
