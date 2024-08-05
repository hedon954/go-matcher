package glicko2

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Player struct {
	entry.Player
	MMR            float64
	Star           int
	startMatchSec  int64
	finishMatchSec int64
	rank           int
}

func CreatePlayer(p entry.Player, info *pto.Glicko2Info) *Player {
	return &Player{
		Player:         p,
		MMR:            info.MMR,
		Star:           info.Star,
		startMatchSec:  info.StartMatchSec,
		finishMatchSec: 0,
		rank:           info.Rank,
	}
}

func (p *Player) GetID() string {
	return p.UID()
}

func (p *Player) IsAi() bool {
	return false
}

func (p *Player) GetMMR() float64 {
	return p.MMR
}

func (p *Player) GetStar() int {
	return p.Star
}

func (p *Player) GetStartMatchTimeSec() int64 {
	return p.startMatchSec
}

func (p *Player) SetStartMatchTimeSec(t int64) {
	p.startMatchSec = t
}

func (p *Player) GetFinishMatchTimeSec() int64 {
	return p.finishMatchSec
}

func (p *Player) SetFinishMatchTimeSec(t int64) {
	p.finishMatchSec = t
}

func (p *Player) GetRank() int {
	return p.rank
}
