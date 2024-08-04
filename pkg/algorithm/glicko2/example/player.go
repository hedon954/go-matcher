package example

import (
	"sync"

	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
	glicko "github.com/zelenin/go-glicko2"
)

type Player struct {
	sync.RWMutex `json:"-"`

	ID      string `json:"id"`
	isAi    bool
	aiLevel int64

	MMR float64 `json:"mmr"`
	RD  float64 `json:"rd"`
	V   float64 `json:"v"`

	rank int
	star int

	startMatchTime  int64
	finishMatchTime int64

	WinWeight int `json:"win_weight"`

	*glicko.Player `json:"-"`
}

func NewPlayer(id string, isAi bool, aiLevel int64, star int, args glicko2.Args) *Player {
	/**
	TODO:
	算法刚启动的时候，会手动配置不同段位的补偿分数，把各个段位的分数人工区分初始化玩家的 MMR，rd 和 v，

		赛季重置也会初始化玩家的相关分数，重置方式如下；
	1. 初始评分（MMR），转换为上赛季评分*70%，最低分1000，最高分7500
	2. RD，根据当前赛季和历史最高赛季的星星数差距决定，最高700，最低为0
	3. 波动率，初始为0.06
	*/
	return &Player{
		RWMutex: sync.RWMutex{},
		ID:      id,
		isAi:    isAi,
		aiLevel: aiLevel,
		MMR:     args.MMR,
		RD:      args.RD,
		V:       args.V,
		star:    star,
		Player:  glicko.NewPlayer(glicko.NewRating(args.MMR, args.RD, args.V)),
	}
}

func (p *Player) GlickoPlayer() *glicko.Player {
	return p.Player
}

func (p *Player) IsAi() bool {
	return p.isAi
}

func (p *Player) GetStartMatchTimeSec() int64 {
	return p.startMatchTime
}

func (p *Player) SetStartMatchTimeSec(t int64) {
	p.startMatchTime = t
}

func (p *Player) GetFinishMatchTimeSec() int64 {
	return p.finishMatchTime
}

func (p *Player) SetFinishMatchTimeSec(t int64) {
	p.finishMatchTime = t
}

func (p *Player) GetID() string {
	return p.ID
}

func (p *Player) GetMMR() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.MMR
}

func (p *Player) GetStar() int {
	return p.star
}

func (p *Player) GetRank() int {
	return p.rank
}

func (p *Player) SetRank(rank int) {
	p.rank = rank
}
