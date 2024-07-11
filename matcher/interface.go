package matcher

import (
	"matcher/common"
	"matcher/pto"
)

type Matcher interface {
	BindRestore(int, string) (int, interface{})
	CreateGroup(*pto.CreateGroup) (int64, error)
	InviteFriend(*InviteFriend) error
	HandleInvite(*HandleInvite) error
	Kick(string, string) error
	ExitGroup(string) error
	DissolveGroup(string) error
	ExitGameGroup(common.Player)
	Match(*pto.PlayerInfo, *MatchReply) error
	CancelMatch(*CancelMatch) error
	MatchSuccess(common.Group)
	UploadAttr(*UploadAttr) error
	ChatInvite(*ChatInvite, *ChatInviteRsp) error
	BroadcastMessage(*RpcChatMessage)
	SetVoiceState(string, int) error
	JoinGroup(*JoinGroup) error
	SyncGroup(group *SyncGroup) error
	SetNearbyJoinGroup(*UserSetting) error
	SetRecentJoinGroup(*UserSetting) error
	RenewGroup(common.Group) []string
}
