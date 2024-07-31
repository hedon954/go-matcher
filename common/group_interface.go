package common

import (
	"github.com/hedon954/go-matcher/pto"
)

// 队伍状态
type GroupState int

const (
	GroupStateDissolve GroupState = 0
	GroupStateInvite   GroupState = 1
	GroupStateQueuing  GroupState = 2
	GroupStateMatched  GroupState = 3
)

// 组队卡片来源
type InviteCardSrc string

const (
	SrcSingleChat  InviteCardSrc = "single"
	SrcClanChat    InviteCardSrc = "clan"
	SrcShare       InviteCardSrc = "share"
	SrcChannelChat InviteCardSrc = "channel"
)

var InviteCardSrcs = []InviteCardSrc{SrcSingleChat, SrcClanChat, SrcShare, SrcChannelChat}

// 聊天卡片状态
type ChatCardState int

const (
	ChatGroupCreate ChatCardState = 1
	ChatGroupUpdate ChatCardState = 2
	ChatGroupExpire ChatCardState = 3
)

type Group interface {
	GroupID() int64
	Base() *GroupBase
	CheckState(validStates ...GroupState) error
	CheckInvite() error
	AddInvitedPlayer(uid string)
	CheckHandleInviteExpired(uid string, srcType pto.InvitationSrcType) error
	IsFull() bool
	GetMatchStrategy() int
	GetGameMode() int
	GetModeVersion() int
	InitUnReadyMap()
}
