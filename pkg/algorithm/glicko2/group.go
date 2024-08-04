package glicko2

// GroupState 队伍状态
type GroupState uint8

const (
	GroupStateUnready GroupState = iota // 未准备
	GroupStateQueuing                   // 匹配中
	GroupStateMatched                   // 匹配完成
)

// GroupType 车队类型
type GroupType uint8

const (
	// 车队类型
	GroupTypeNotTeam GroupType = iota
	GroupTypeNormalTeam
	GroupTypeUnfriendlyTeam
	GroupTypeMaliciousTeam
)

// Group 是一个队伍，
// 玩家可以自行组队，单个玩家开始匹配的时候也会为其单独创建一个队伍，
// 匹配前后队伍都不会被拆开。
type Group interface {

	// 队伍 ID
	GetID() string

	// 获取队伍里的玩家列表
	GetPlayers() []Player

	// 队伍中的玩家个数
	PlayerCount() int

	// 获取队伍 mmr 值
	GetMMR() float64

	// 获取队伍段位值
	GetStar() int

	// 队伍状态
	GetState() GroupState
	SetState(state GroupState)

	// 开始匹配的时间，取 player 中最早的
	GetStartMatchTimeSec() int64
	SetStartMatchTimeSec(t int64)

	// 结束匹配的时间
	GetFinishMatchTimeSec() int64
	SetFinishMatchTimeSec(t int64)

	// 获取车队类型
	Type() GroupType

	// 当返回 true 时，会自动填充 is_ai 组成房间，第二个返回值为填充的 ai team
	CanFillAi() bool

	// 强制退出时对每个玩家的处理逻辑
	ForceCancelMatch(reason string, waitSec int64)

	// IsNewer 判断组队是否认定为新手
	IsNewer() bool
}
