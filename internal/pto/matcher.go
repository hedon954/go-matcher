package pto

import (
	"github.com/hedon954/go-matcher/internal/constant"
)

// PlayerInfo defines the common information of a player.
// It is always used to initial a player.
type PlayerInfo struct {
	UID           string                 `json:"uid" binding:"required"`
	GameMode      constant.GameMode      `json:"game_mode" binding:"required"`
	ModeVersion   int64                  `json:"mode_version" binding:"required"`
	MatchStrategy constant.MatchStrategy `json:"match_strategy" binding:"required"`

	Glicko2Info *Glicko2Info `json:"glicko2_info"`
}

type EnterGroup struct {
	PlayerInfo
	Source EnterGroupSourceType
}

type CreateGroup struct {
	PlayerInfo
}

type GroupUser struct {
	GroupID     int64
	Owner       string
	UidList     []string
	GameMode    int
	RingList    []int
	ModeVersion int
	MaxStar     int
	PropList    []int
	StateList   []int
	RoleList    []int
	SkinList    []int
	SuitList    []int
	ShowList    []int
}

// EnterGroupSourceType is the source type of entering a group.
type EnterGroupSourceType int

const (
	EnterGroupSourceTypeInvite                            = 0 // invited by other
	EnterGroupSourceTypeNearby       EnterGroupSourceType = 1 // from recent list
	EnterGroupSourceTypeRecent       EnterGroupSourceType = 2 // from nearby list
	EnterGroupSourceTypeFriend       EnterGroupSourceType = 3 // from friend list
	EnterGroupSourceTypeWorldChannel EnterGroupSourceType = 4 // from world channel
	EnterGroupSourceTypeClanChannel  EnterGroupSourceType = 5 // from clan channel
	EnterGroupSourceTypeShare        EnterGroupSourceType = 6 // from share link
)

type InviteMsg struct {
	InviterUID  string
	InviteeUID  string
	Source      EnterGroupSourceType
	GameMode    constant.GameMode
	ModeVersion int64
}

type HandleInvite struct {
	InviteUid  string // 谁邀请你的
	HandleType int
	SrcType    EnterGroupSourceType
	Player     *PlayerInfo // 当前玩家
	Message    string
	Platform   int
}

type UploadAttr struct {
	Uid string
	Attribute
	// Extra 是每个玩法特定的数据，需要在玩法内部具体处理
	Extra []byte
}

type Attribute struct {
	SkinID     int
	Rank       int
	SkinSkills []*SkinSkill
	KillEffect int
	Platform   int
	SkinIdList []int
	ShowLevel  int
	ShowSkill  string
	ShowValue  int
	Nickname   string
}

type SkinSkill struct {
	SkillType int     `json:"skill_type"`
	Value     float64 `json:"value"`
}

type CheckInviteFriend struct {
	Platform    int
	GameMode    int
	ModeVersion int
	InviteUid   string
	FriendUid   string
	Source      EnterGroupSourceType
	NoCheck     bool
}

type UserVoiceState struct {
	UID   string
	State int
}

type MatchInfo struct {
}
