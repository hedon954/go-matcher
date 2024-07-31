package pto

type CreateGroup struct {
	UID               string
	GameMode          int
	ModeVersion       int
	MatchStrategy     int
	UnityNamespacePre string
}

type GroupUser struct {
	GroupId     int
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

// 邀请渠道
type InvitationSrcType int

const (
	InvitationSrcSingleChat InvitationSrcType = 1  // 组队页面(好友tab)
	InvitationSrcClanRank   InvitationSrcType = 2  // 排位赛(聊天)
	InvitationSrcClanRace   InvitationSrcType = 3  // 战队赛(聊天)
	InvitationSrcNearBy     InvitationSrcType = 4  // 附近的人
	InvitationSrcShare      InvitationSrcType = 5  // 应用外邀请
	InvitationSrcRace       InvitationSrcType = 7  // 公开赛
	InvitationSrcRecent     InvitationSrcType = 8  // 最近tab
	InvitationSrcClan       InvitationSrcType = 9  // 战队tab
	InvitationSrcChannel    InvitationSrcType = 10 // 世界频道
)

// 好友游戏邀请处理类型：接受或拒绝
const (
	InviteHandleTypeAccept    = 1
	InviteHandleTypeRefuse    = 2
	InviteHandleTypeRefuseMsg = 3
	DefaultInviteRefuseMsg    = "不好意思，现在不方便，下次约。"
)

type InviteFriend struct {
	InviteUid string
	FriendUid string
	Source    InvitationSrcType
	Platform  int
	NoCheck   bool
}

type HandleInvite struct {
	InviteUid  string // 谁邀请你的
	HandleType int
	SrcType    InvitationSrcType
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
	Source      InvitationSrcType
	NoCheck     bool
}
