package goatgameglicko2

import (
	"github.com/hedon954/go-matcher/internal/entry"
	glicko "github.com/zelenin/go-glicko2"
)

type Player struct {
	entry.Player
	MMR            float64
	Star           int
	glickoPlayer   glicko.Player
	startMatchSec  int64
	finishMatchSec int64
	rank           int
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

func (p *Player) GlickoPlayer() *glicko.Player {
	return &p.glickoPlayer
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
