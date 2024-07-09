package common

import (
	"sync"

	"matcher/pto"
)

// Base 是 Player 的基础类，所有游戏模式和所有匹配策略共用
type Base struct {
	// 在 Base 的内部方法中不要进行同步处理，统一交给外部方法调用
	sync.RWMutex
	uid               string
	GroupID           int64
	onlineState       OnlineState
	voiceState        VoiceState
	WithCouple        bool
	GameMode          int
	ModeVersion       int
	MatchStrategy     int
	UnityNamespace    string
	UnityNamespacePre string
}

func NewBase(pInfo *pto.PlayerInfo) *Base {
	return &Base{
		uid:               pInfo.UID,
		onlineState:       OnlineStateOnline,
		voiceState:        VoiceStateOff,
		GameMode:          pInfo.GameMode,
		MatchStrategy:     pInfo.MatchStrategy,
		ModeVersion:       pInfo.ModeVersion,
		UnityNamespacePre: pInfo.UnityNamespacePre,
	}
}

func (b *Base) Inner() *Base {
	return b
}

func (b *Base) UID() string {
	return b.uid
}

func (b *Base) GetOnlineState() OnlineState {
	return b.onlineState
}

func (b *Base) SetOnlineState(state OnlineState) {
	b.onlineState = state
}
