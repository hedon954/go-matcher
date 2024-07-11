package matcher

import (
	"matcher/pto"
)

type InviteFriend struct {
	InviteUid string
	FriendUid string
	Source    pto.InvitationSrcType
	Platform  int
	NoCheck   bool
}

type CheckInviteFriend struct {
	Platform    int
	GameMode    int
	ModeVersion int
	InviteUid   string
	FriendUid   string
	Source      pto.InvitationSrcType
	NoCheck     bool
}

type HandleInvite struct {
	InviteUid  string // 谁邀请你的
	HandleType int
	SrcType    pto.InvitationSrcType
	Player     *pto.PlayerInfo // 当前玩家
	Message    string
	Platform   int
}

type MatchReply struct {
	GroupID  int
	Estimate int
	NeedNum  int
}

type CancelMatch struct {
	Uid        string
	Platform   int
	RoomMember int
}

type UploadAttr struct {
	Uid             string
	SkinID          int
	Rank            int
	SkinSkills      []*SkinSkill
	KillEffect      int
	Platform        int
	SkinIdList      []int
	ShowLevel       int
	ShowSkill       string
	ShowValue       int
	QualifyingLevel int
	Nickname        string
	AccPoint        int
	RatedLevel      int
	BountyTitle     string
	PropId          int
	ClothesId       int
	SuitId          int
	TeamBroadcastId int
	TeamBubbleId    int
	KuromiInfo      *KuromiInfo
	GoatInfo        *GoatInfo
}

type SkinSkill struct {
	SkillType int     `json:"skill_type"`
	Value     float64 `json:"value"`
}

type KuromiInfo struct {
	TotalPvpCount int
	TodayPvpCount int
}

type GoatInfo struct {
	TotalPvpCount int
	TodayPvpCount int
}

type ChatInvite struct {
	InviteUid string
	FriendUid string
	ChatType  int
	Platform  int
}

type ChatInviteRsp struct {
	GroupMemberCount int
}

type RpcChatMessage struct { // json存入redis,应同ChatMessage
	Nickname     string `json:"-"`
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"` // 自己发的邀请消息会存入自己的列表，用于区分对象
	Content      string `json:"content"`
	ChatType     int    `json:"chat_type"`
	MessageType  int    `json:"message_type"`
	ChatId       int    `json:"chat_id"` // msgNo，世界频道是channelId
	Timestamp    int    `json:"timestamp"`
	Ultimate     int    `json:"-"`
	Star         int    `json:"-"`
	IsSent       bool   `json:"-"`
	Ext          string `json:"ext"`
	Platform     int    `json:"-"`
	PushMessage  string `json:"-"`
	Charm        int    `json:"charm"`
	RatedLevel   int    `json:"-"`
	Market       string `json:"-"`
	SubType      int    `json:"sub_type"`
	DeviceId     string `json:"-"`
	RiskLevel    int    `json:"risk_level"`
	RiskType     int    `json:"-"`
	RiskTypeDesc string `json:"-"`
	Shumei       string `json:"-"`
	MessageID    string `json:"-"` // 世界频道消息标识，其他场景未使用
}

type JoinGroup struct {
	Player       *pto.PlayerInfo
	GroupUserUid string
	Ultimate     int
	Star         int
	Charm        int
	Platform     int
	Source       int
}

type SyncGroup struct {
	Uid       string
	GroupID   int
	GroupName string
}
type UserSetting struct {
	Uid      string
	Setting  bool
	Platform int
}
