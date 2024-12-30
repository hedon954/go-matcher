package entry

import (
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/concurrent"
)

// Player represents a player in a Group.
type Player interface {
	Coder
	Base() *PlayerBase
	UID() string
	GetPlayerInfo() *pto.PlayerInfo
	SetAttr(attr *pto.UploadPlayerAttr) error
}

// PlayerOnlineState is the state of a player.
// TODO: try to use state machine to manage player state.
type PlayerOnlineState int8

const (
	PlayerOnlineStateOffline  PlayerOnlineState = 0
	PlayerOnlineStateOnline   PlayerOnlineState = 1
	PlayerOnlineStateInGroup  PlayerOnlineState = 2
	PlayerOnlineStateInMatch  PlayerOnlineState = 3
	PlayerOnlineStateInGame   PlayerOnlineState = 4
	PlayerOnlineStateInSettle PlayerOnlineState = 5
)

// PlayerVoiceState is the voice state of a player.
type PlayerVoiceState int8

const (
	PlayerVoiceStateMute   PlayerVoiceState = 0
	PlayerVoiceStateUnmute PlayerVoiceState = 1
)

// PlayerBase holds the common fields of a Player for all kinds of game mode and match strategy.
type PlayerBase struct {
	// ReentrantLock is a reentrant lock support multiple locks in the same goroutine
	// Use it to help avoid deadlock.
	lock          *concurrent.ReentrantLock
	uid           string
	IsAI          bool
	GroupID       int64
	matchStrategy constant.MatchStrategy

	onlineState PlayerOnlineState
	voiceState  PlayerVoiceState

	// TODO: other common attributes
	pto.PlayerInfo
	pto.Attribute
}

func NewPlayerBase(info *pto.PlayerInfo) *PlayerBase {
	b := &PlayerBase{
		lock:        new(concurrent.ReentrantLock),
		uid:         info.UID,
		onlineState: PlayerOnlineStateOnline,
		voiceState:  PlayerVoiceStateMute,
		PlayerInfo:  *info,
	}

	return b
}

func (p *PlayerBase) GetPlayerInfo() *pto.PlayerInfo {
	return &p.PlayerInfo
}

func (p *PlayerBase) Base() *PlayerBase {
	return p
}

func (p *PlayerBase) UID() string {
	return p.uid
}

// CheckOnlineState checks if the player is in a valid online state.
func (p *PlayerBase) CheckOnlineState(valids ...PlayerOnlineState) error {
	for _, vs := range valids {
		if p.onlineState == vs {
			return nil
		}
	}
	switch p.onlineState {
	case PlayerOnlineStateOffline:
		return merr.ErrPlayerOffline
	case PlayerOnlineStateOnline:
		return merr.ErrPlayerNotInGroup
	case PlayerOnlineStateInGroup:
		return merr.ErrPlayerInGroup
	case PlayerOnlineStateInGame:
		return merr.ErrPlayerInGame
	case PlayerOnlineStateInMatch:
		return merr.ErrPlayerInMatch
	case PlayerOnlineStateInSettle:
		return merr.ErrPlayerInSettle
	}
	panic("unreachable")
}

func (p *PlayerBase) SetOnlineState(s PlayerOnlineState) {
	p.onlineState = s
}

func (p *PlayerBase) GetOnlineState() PlayerOnlineState {
	return p.onlineState
}

func (p *PlayerBase) SetOnlineStateWithLock(s PlayerOnlineState) {
	p.Lock()
	defer p.Unlock()
	p.onlineState = s
}

func (p *PlayerBase) GetOnlineStateWithLock() PlayerOnlineState {
	p.Lock()
	defer p.Unlock()
	return p.onlineState
}

func (p *PlayerBase) SetVoiceState(s PlayerVoiceState) {
	p.voiceState = s
}

func (p *PlayerBase) GetVoiceState() PlayerVoiceState {
	return p.voiceState
}

func (p *PlayerBase) SetAttr(attr *pto.UploadPlayerAttr) error {
	p.Attribute = attr.Attribute
	return nil
}

func (p *PlayerBase) GetMatchStrategy() constant.MatchStrategy {
	return p.matchStrategy
}

func (p *PlayerBase) SetMatchStrategy(s constant.MatchStrategy) {
	p.matchStrategy = s
}

func (p *PlayerBase) GetMatchStrategyWithLock() constant.MatchStrategy {
	p.Lock()
	defer p.Unlock()
	return p.matchStrategy
}

func (p *PlayerBase) SetMatchStrategyWithLock(s constant.MatchStrategy) {
	p.Lock()
	defer p.Unlock()
	p.matchStrategy = s
}

func (p *PlayerBase) Lock() {
	p.lock.Lock()
}
func (p *PlayerBase) Unlock() {
	p.lock.Unlock()
}
