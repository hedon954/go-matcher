package glicko2

// Player 是一个玩家的抽象
type Player interface {

	// 玩家ID
	GetID() string

	// 是否是 AI
	IsAi() bool

	// 获取 mmr 值
	GetMMR() float64

	// 段位值
	GetStar() int

	// 开始匹配的时间
	GetStartMatchTimeSec() int64
	SetStartMatchTimeSec(t int64)

	// 结束匹配的时间
	GetFinishMatchTimeSec() int64
	SetFinishMatchTimeSec(t int64)

	// 赛后在阵营内的排名
	GetRank() int
}
