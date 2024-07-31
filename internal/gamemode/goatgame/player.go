package goatgame

import (
	"fmt"

	"github.com/hedon954/go-matcher/common"
	"github.com/hedon954/go-matcher/enum"
	"github.com/hedon954/go-matcher/internal/gamemode/goatgame/glicko2"
)

func CreatePlayer(base *common.PlayerBase) (common.Player, error) {
	switch base.MatchStrategy {
	case enum.MatchStrategyGlicko2:
		return glicko2.NewPlayer(base)
	default:
		// TODO: or select a default match strategy?
		return nil, fmt.Errorf("unknown match strategy: %d", base.MatchStrategy)
	}
}
