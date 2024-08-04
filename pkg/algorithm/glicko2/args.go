package glicko2

// Args 封装了 glicko-2 算法的 3 个核心参数
type Args struct {

	// 用户游戏评分，这是对玩家能力的直接衡量。
	MMR float64 `json:"mmr"`

	// 评分偏差，这是对评分准确性的衡量。
	// 如果你是一个新玩家或者很久没玩游戏了，你的 RD 会很高，表示你的真实技能可能和你的评分相差很大。
	// 如果你经常玩游戏，你的 RD 会降低，表示你的评分越来越接近你的真实技能。
	RD float64 `json:"rd"`

	// 波动率；这是对玩家的评分变动幅度的衡量。
	V float64 `json:"v"`
}
