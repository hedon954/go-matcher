package common

import (
	"github.com/hedon954/go-matcher/pto"
)

type PlayerOnlineState int

const (
	PlayerOnlineStateOffline PlayerOnlineState = 1
	PlayerOnlineStateOnline  PlayerOnlineState = 2
	PlayerOnlineStateGroup   PlayerOnlineState = 3
	PlayerOnlineStateQueuing PlayerOnlineState = 4
	PlayerOnlineStateGame    PlayerOnlineState = 5
	PlayerOnlineStateSettle  PlayerOnlineState = 6
)

type PlayerVoiceState int

const (
	PlayerVoiceStateOff PlayerVoiceState = 0
	PlayerVoiceStateOn  PlayerVoiceState = 1
)

type Player interface {
	// Base 是 Player 的内部类，用于定义所有模式所有玩法通用的基类。
	//
	// 这里为什么这么设计？
	//  考虑到通用的结构基本比较固定，如果纯搞接口的话，需要为每个字段搞一套 getter 和 setter，会比较繁琐。
	//  这里引入的耦合性，是可以接受的，因为如果这些通用字段都发生了变化的话，那基本上整个匹配服架构都发生了变化。
	//  所以这里妥协了耦合性，带来了比较多的便利性。
	Base() *PlayerBase

	// UID 获取玩家的 uid
	UID() string

	// VersionMatched 判断玩家版本是否一致，只有版本一致的玩家才可以在一起玩
	VersionMatched(Player) bool
	SetWithCouple(b bool)
	SetAttr(attr *pto.UploadAttr) error
}

type OnlinePlayer interface {
	Player
	CreateGroup() (GroupPlayer, error)
	JoinGroup() (GroupPlayer, error)
}

type GroupPlayer interface {
	Player
	ExitGroup() (OnlinePlayer, error)
	DissolveGroup() (OnlinePlayer, error)
	Ready() (ReadyPlayer, error)
	StartMatch() (QueuePlayer, error)
}

type ReadyPlayer interface {
	Player
	UnReady() (GroupPlayer, error)
	StartMatch() (QueuePlayer, error)
}

type QueuePlayer interface {
	Player
	CancelMatch() (ReadyPlayer, error)
	MatchTimeout() (ReadyPlayer, error)
	MatchSuccess() (GamingPlayer, error)
}

type GamingPlayer interface {
	Player
	GameEnd() (SettlingPlayer, error)
	Escape() (OnlinePlayer, error)
}

type SettlingPlayer interface {
	Player
	Finish() (GroupPlayer, error)
}
