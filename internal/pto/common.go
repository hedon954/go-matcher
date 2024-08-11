package pto

import (
	"github.com/hedon954/go-matcher/internal/constant"
)

// PlayerInfo defines the common information of a player.
// It is often used to initial a player.
type PlayerInfo struct {
	UID         string            `json:"uid" binding:"required"`
	GameMode    constant.GameMode `json:"game_mode" binding:"required" example:"905"`
	ModeVersion int64             `json:"mode_version" binding:"required" example:"1"`
	Star        int64             `json:"star"`
	Rank        int64             `json:"rank"`

	Glicko2Info *Glicko2Info `json:"glicko2_info"`
}

// GameResult defines the common information of a game result.
type GameResult struct {
	RoomID         int64
	StartTime      int64
	EndTime        int64
	ScoreSign      string
	GameMode       constant.GameMode
	MatchStrategy  constant.MatchStrategy
	ModeVersion    int64
	PlayerState    map[string]GamePlayerState
	AIPlayer       map[string]bool
	PlayerMetaInfo map[string]PlayerMetaInfo

	// Result is the game result detail info.
	// It is sent from client and different game mode would be different.
	Result []byte
}

type PlayerMetaInfo struct {
	/* add according to requirements */
}

type GamePlayerState uint8

const (
	Offline GamePlayerState = 0
	Online  GamePlayerState = 1
)
