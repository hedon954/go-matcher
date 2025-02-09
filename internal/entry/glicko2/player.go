package glicko2

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type PlayerBaseGlicko2 struct {
	*entry.PlayerBase
	MMR            float64
	Star           int64
	StartMatchSec  int64
	FinishMatchSec int64
	Rank           int64
}

func CreatePlayerBase(p *entry.PlayerBase, info *pto.Glicko2Info) *PlayerBaseGlicko2 {
	return &PlayerBaseGlicko2{
		PlayerBase: p,
		MMR:        info.MMR,
		Star:       info.Star,
		Rank:       info.Rank,
	}
}

func (p *PlayerBaseGlicko2) GetID() string {
	return p.UID()
}

func (p *PlayerBaseGlicko2) IsAi() bool {
	return false
}

func (p *PlayerBaseGlicko2) GetMMR() float64 {
	return p.MMR
}

func (p *PlayerBaseGlicko2) GetStar() int {
	return int(p.Star)
}

func (p *PlayerBaseGlicko2) GetStartMatchTimeSec() int64 {
	return p.StartMatchSec
}

func (p *PlayerBaseGlicko2) SetStartMatchTimeSec(t int64) {
	p.StartMatchSec = t
}

func (p *PlayerBaseGlicko2) GetFinishMatchTimeSec() int64 {
	return p.FinishMatchSec
}

func (p *PlayerBaseGlicko2) SetFinishMatchTimeSec(t int64) {
	p.FinishMatchSec = t
}

func (p *PlayerBaseGlicko2) GetRank() int {
	return int(p.Rank)
}
