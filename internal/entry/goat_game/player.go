package goat_game

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Player struct {
	*glicko2.PlayerBaseGlicko2
}

func CreatePlayer(base *entry.PlayerBase, pInfo *pto.PlayerInfo) entry.Player {
	p := &Player{}
	// ... other common fields

	p.withMatchStrategy(base, pInfo.Glicko2Info)
	return p
}

func (p *Player) withMatchStrategy(base *entry.PlayerBase, info *pto.Glicko2Info) {
	p.PlayerBaseGlicko2 = glicko2.CreatePlayerBase(base, info)
	// ... other match strategy initialization
}
