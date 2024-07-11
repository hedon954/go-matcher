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
