package glicko2

// Room 是一个房间的抽象，由多个 team 组成
type Room interface {

	// 房间id
	GetID() int64

	// 获取玩家中的阵营
	GetTeams() []Team

	// 通过结果排名获取阵营列表
	SortTeamByRank() []Team

	// 向房间中添加阵营
	AddTeam(t Team)

	// 向房间中移除阵营
	RemoveTeam(t Team)

	// 获取房间的 mmr 值
	GetMMR() float64

	// 房间的玩家个数
	PlayerCount() int

	// 房间的开始匹配时间，取最早的玩家的
	GetStartMatchTimeSec() int64

	// 房间的完成匹配时间
	GetFinishMatchTimeSec() int64
	SetFinishMatchTimeSec(t int64)

	// 是否存在 ai
	HasAi() bool
}
