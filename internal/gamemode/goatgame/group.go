package goatgame

import (
	"fmt"

	"github.com/hedon954/go-matcher/common"
	"github.com/hedon954/go-matcher/enum"
	"github.com/hedon954/go-matcher/internal/gamemode/goatgame/glicko2"
)

func CreateGroup(groupID int64, p common.Player) (common.Group, error) {
	switch p.Base().MatchStrategy {
	case enum.MatchStrategyGlicko2:
		return glicko2.NewGroup(groupID, p)
	default:
		// TODO: or select a default match strategy?
		return nil, fmt.Errorf("unknown match strategy: %d", p.Base().MatchStrategy)
	}
}
