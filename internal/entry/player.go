package entry

import (
	"sync"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
)

// Player represents a player in a Group.
type Player interface {
	Base() *PlayerBase
	UID() string
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
	sync.RWMutex
	uid           string
	GroupID       int64
	onlineState   PlayerOnlineState
	VoiceState    PlayerVoiceState
	GameMode      constant.GameMode
	ModeVersion   int64
	MatchStrategy constant.MatchStrategy

	// TODO: other common attributes
}

func NewPlayerBase(info *pto.PlayerInfo) *PlayerBase {
	b := &PlayerBase{
		uid:           info.UID,
		onlineState:   PlayerOnlineStateOnline,
		VoiceState:    PlayerVoiceStateMute,
		GameMode:      info.GameMode,
		ModeVersion:   info.ModeVersion,
		MatchStrategy: info.MatchStrategy,
	}

	return b
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

func (p *PlayerBase) SetVoiceState(s PlayerVoiceState) {
	p.VoiceState = s
}
func (p *PlayerBase) GetVoiceState() PlayerVoiceState {
	return p.VoiceState
}
