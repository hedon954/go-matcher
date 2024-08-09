package pto

import (
	"github.com/hedon954/go-matcher/internal/constant"
)

// PlayerInfo defines the common information of a player.
// It is always used to initial a player.
type PlayerInfo struct {
	UID         string            `json:"uid" binding:"required"`
	GameMode    constant.GameMode `json:"game_mode" binding:"required"`
	ModeVersion int64             `json:"mode_version" binding:"required"`

	Glicko2Info *Glicko2Info `json:"glicko2_info"`
}
