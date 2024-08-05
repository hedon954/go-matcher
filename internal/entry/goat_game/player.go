package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game/glicko2"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Player struct {
	*entry.PlayerBase
}

func CreatePlayer(base *entry.PlayerBase, pInfo *pto.PlayerInfo) (entry.Player, error) {
	player := &Player{PlayerBase: base}
	switch base.MatchStrategy {
	case constant.MatchStrategyGlicko2:
		return glicko2.CreatePlayer(player, pInfo.Glicko2Info), nil
	default:
		return nil, fmt.Errorf("unknown match strategy: %d", base.MatchStrategy)
	}
}
