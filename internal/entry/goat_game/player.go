package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Player struct {
	// combined with the struct implementing the matching strategy
	// to support different matching strategies
	*glicko2.PlayerBaseGlicko2

	// some game mode specific fields
	TotalPvpCount int64
	TodayPvpCount int64
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

// SetAttr rewrite the base method if needed.
func (p *Player) SetAttr(attr *pto.UploadPlayerAttr) error {
	if err := p.Base().SetAttr(attr); err != nil {
		return err
	}

	attribute, err := pto.FromAttrJson[Attribute](attr)
	if err != nil {
		return fmt.Errorf("invalid attribute: %w", err)
	}

	p.Glicko2Info.MMR = attribute.MMR
	p.Glicko2Info.Rank = attribute.Rank
	p.Glicko2Info.Star = attribute.Star

	p.TotalPvpCount = attribute.TotalPvpCount
	p.TodayPvpCount = attribute.TodayPvpCount
	return nil
}

type Attribute struct {
	MMR           float64 `json:"mmr"`
	Star          int     `json:"star"`
	Rank          int     `json:"rank"`
	TotalPvpCount int64   `json:"total_pvp_count"`
	TodayPvpCount int64   `json:"today_pvp_count"`
}
