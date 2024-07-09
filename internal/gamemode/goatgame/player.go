package goatgame

import (
	"fmt"

	"matcher/common"
	"matcher/enum"
	"matcher/internal/gamemode/goatgame/glicko2"
)

func CreatePlayer(base *common.Base) (common.Player, error) {
	switch base.MatchStrategy {
	case enum.MatchStrategyGlicko2:
		return glicko2.NewPlayer(base)
	default:
		// TODO: or select a default match strategy?
		return nil, fmt.Errorf("unknown match strategy: %d", base.MatchStrategy)
	}
}
