package matcher

import (
	"github.com/hedon954/go-matcher/common"
	"github.com/hedon954/go-matcher/pto"
)

// Matcher 定义了匹配器的通用外围接口
type Matcher interface {
	BindRestore(int, string) (state int, matchInfo interface{})
	CreateGroup(info *pto.PlayerInfo) (int64, error)
	InviteFriend(*pto.InviteFriend) error
	HandleInvite(*pto.HandleInvite) error
	Kick(string, string) error
	ExitGroup(string) error
	DissolveGroup(string) error
	Match(*pto.PlayerInfo, *MatchReply) error
	CancelMatch(*CancelMatch) error
	UploadAttr(*pto.UploadAttr) error
	ChatInvite(*ChatInvite, *ChatInviteRsp) error
	BroadcastMessage(*RpcChatMessage)
	SetVoiceState(string, common.PlayerVoiceState) error
	JoinGroup(*JoinGroup) error
	SyncGroup(group *SyncGroup) error
	SetNearbyJoinGroup(*UserSetting) error
	SetRecentJoinGroup(*UserSetting) error
	RenewGroup(common.Group) []string
	GroupReady(string, int, int) error
	ChangeRole(string, int, int) error
}

// GameHelper 定义了游戏辅助器的通用接口，主要是一些跟游戏模式相关的判断逻辑
type GameHelper interface {
	IsBanGame(uid string) bool
	IsInGameTime(nowUnix int64) bool
}
