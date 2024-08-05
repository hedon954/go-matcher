package glicko2

// Team 是一个阵营的抽象，由 1~n 个 Group 组成
type Team interface {

	// 获取 groups 列表
	GetGroups() []Group

	// 添加 group 到 team 中
	AddGroup(group Group)

	// 移除 group
	RemoveGroup(groupId string)

	// 玩家数量
	PlayerCount() int

	// 阵营的 mmr
	GetMMR() float64

	// 阵营的段位值
	GetStar() int

	// 阵营的开始匹配时间，取最早的玩家的
	GetStartMatchTimeSec() int64

	// 阵营的完成匹配时间
	GetFinishMatchTimeSec() int64
	SetFinishMatchTimeSec(t int64)

	// 当前阵营是否是 ai
	IsAi() bool

	// 是否可以填充 AI
	CanFillAi() bool

	// Team是否已满
	IsFull(teamPlayerLimit int) bool

	// IsNewer 判断 Team 是否认定为新手
	IsNewer() bool
}
