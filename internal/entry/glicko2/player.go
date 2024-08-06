package glicko2

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type PlayerBaseGlicko2 struct {
	*entry.PlayerBase
	MMR            float64
	Star           int
	startMatchSec  int64
	finishMatchSec int64
	rank           int
}

func CreatePlayerBase(p *entry.PlayerBase, info *pto.Glicko2Info) *PlayerBaseGlicko2 {
	return &PlayerBaseGlicko2{
		PlayerBase:     p,
		MMR:            info.MMR,
		Star:           info.Star,
		startMatchSec:  info.StartMatchSec,
		finishMatchSec: 0,
		rank:           info.Rank,
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
	return p.Star
}

func (p *PlayerBaseGlicko2) GetStartMatchTimeSec() int64 {
	return p.startMatchSec
}

func (p *PlayerBaseGlicko2) SetStartMatchTimeSec(t int64) {
	p.startMatchSec = t
}

func (p *PlayerBaseGlicko2) GetFinishMatchTimeSec() int64 {
	return p.finishMatchSec
}

func (p *PlayerBaseGlicko2) SetFinishMatchTimeSec(t int64) {
	p.finishMatchSec = t
}

func (p *PlayerBaseGlicko2) GetRank() int {
	return p.rank
}
