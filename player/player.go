package player

type OnlineState int

const (
	OnlineStateOffline OnlineState = 1
	OnlineStateOnline  OnlineState = 2
	OnlineStateGroup   OnlineState = 3
	OnlineStateQueue   OnlineState = 4
	OnlineStateGame    OnlineState = 5
	OnlineStateSettle  OnlineState = 6
)

type VoiceState int

const (
	VoiceStateOff VoiceState = 0
	VoiceStateOn  VoiceState = 1
)

type Player interface {
	// Inner 是 Player 的内部类，用于定义所有模式所有玩法通用的基类。
	// 这里不用 Base() 是为了防止组合的时候字段和方法命名冲突。
	Inner() *Base
	UID() string
}

type OnlinePlayer interface {
	CreateGroup() (GroupPlayer, error)
	JoinGroup() (GroupPlayer, error)
}

type GroupPlayer interface {
	ExitGroup() (OnlinePlayer, error)
	DissolveGroup() (OnlinePlayer, error)
	Ready() (ReadyPlayer, error)
	StartMatch() (QueuePlayer, error)
}

type ReadyPlayer interface {
	UnReady() (GroupPlayer, error)
	StartMatch() (QueuePlayer, error)
}

type QueuePlayer interface {
	CancelMatch() (ReadyPlayer, error)
	MatchTimeout() (ReadyPlayer, error)
	MatchSuccess() (GamingPlayer, error)
}

type GamingPlayer interface {
	GameEnd() (SettlingPlayer, error)
	Escape() (OnlinePlayer, error)
}

type SettlingPlayer interface {
	Finish() (GroupPlayer, error)
}
