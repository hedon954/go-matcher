package common

import (
	"sync"

	"github.com/hedon954/go-matcher/pto"
)

// PlayerBase 是 Player 的基础类，所有游戏模式和所有匹配策略共用
type PlayerBase struct {
	// 在 PlayerBase 的内部方法中不要进行同步处理，统一交给外部方法调用
	sync.RWMutex
	Uid               string
	Platform          int
	GroupID           int64
	onlineState       PlayerOnlineState
	voiceState        PlayerVoiceState
	WithCouple        bool
	GameMode          int
	ModeVersion       int
	MatchStrategy     int
	UnityNamespace    string
	UnityNamespacePre string

	Attributes pto.Attribute
}

func NewPlayerBase(pInfo *pto.PlayerInfo) *PlayerBase {
	return &PlayerBase{
		Uid:               pInfo.Uid,
		onlineState:       PlayerOnlineStateOnline,
		voiceState:        PlayerVoiceStateOff,
		GameMode:          pInfo.GameMode,
		MatchStrategy:     pInfo.MatchStrategy,
		ModeVersion:       pInfo.ModeVersion,
		UnityNamespacePre: pInfo.UnityNamespacePre,
	}
}

func (b *PlayerBase) Base() *PlayerBase {
	return b
}

func (b *PlayerBase) UID() string {
	return b.Uid
}

func (b *PlayerBase) GetOnlineState() PlayerOnlineState {
	return b.onlineState
}

func (b *PlayerBase) SetOnlineState(state PlayerOnlineState) {
	b.onlineState = state
}

func (b *PlayerBase) GetVoiceState() PlayerVoiceState {
	return b.voiceState
}

func (b *PlayerBase) SetVoiceState(state PlayerVoiceState) {
	b.voiceState = state
}

func (b *PlayerBase) VersionMatched(b2 Player) bool {
	return b.GameMode == b2.Base().GameMode && b.ModeVersion == b2.Base().ModeVersion
}

func (b *PlayerBase) SetWithCouple(with bool) {
	b.WithCouple = with
}

func (b *PlayerBase) SetAttr(attr *pto.UploadAttr) error {
	b.Attributes = attr.Attribute
	return nil
}
