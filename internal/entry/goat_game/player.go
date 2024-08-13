package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
	"github.com/hedon954/go-matcher/internal/pb"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/typeconv"
)

type Player struct {
	// combined with the struct implementing the matching strategy
	// to support different matching strategies
	*glicko2.PlayerBaseGlicko2

	// some game mode specific fields
	TotalPvpCount int64
	TodayPvpCount int64
}

func CreatePlayer(base *entry.PlayerBase, pInfo *pto.PlayerInfo) (entry.Player, error) {
	p := &Player{}
	// ... other common fields

	if pInfo.Glicko2Info == nil {
		return nil, fmt.Errorf("game[%d] need glicko2 info", base.GameMode)
	}
	p.withMatchStrategy(base, pInfo.Glicko2Info)
	return p, nil
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

	if p.Base().GetMatchStrategy() == constant.MatchStrategyGlicko2 {
		if err := p.setGlicko2Attr(attr.Extra); err != nil {
			return err
		}
	}

	return nil
}

func (p *Player) setGlicko2Attr(extra []byte) error {
	attribute, err := typeconv.FromProto[pb.GoatGameAttribute](extra)
	if err != nil {
		return fmt.Errorf("invalid glicko2 attribute: %w", err)
	}
	p.Glicko2Info.MMR = attribute.Mmr
	p.Glicko2Info.Rank = p.Rank
	p.Glicko2Info.Star = p.Star
	return nil
}
