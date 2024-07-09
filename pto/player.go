package pto

type PlayerInfo struct {
	// 确定现在在玩哪个版本的哪个游戏
	GameMode      int
	ModeVersion   int
	MatchStrategy int

	UID           string
	Platform      int
	ClientVersion string
	HotVersion    string

	RobCoinLevel           int
	RobCoinLen             int
	IsRobCoinAi            bool
	Ping                   map[string]int
	NearbyJoinGroupAllowed bool
	RecentJoinGroupAllowed bool
	Market                 string
	GrandBattleVersion     int
	UserRankList           []*UserRank
	SeasonPoints           int
	MatchLevel             int
	SkinID                 int
	HistoryMaxStar         int
	IsNewer                bool
	IsBacker               bool
	GroupName              string
	UnityNamespacePre      string
}

type UserRank struct {
	UID  string
	Rank int
}
