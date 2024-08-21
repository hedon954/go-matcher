package config

import (
	"time"
)

// TODO: support different delay timer type for different game mode
// type GameModeConfig interface {
//    GetDelayTimerConfig() DelayTimerConfig
//    GetMatchStrategy() MatchStrategy
// }

// DelayTimerConfig defines the delay timer config.
type DelayTimerConfig struct {
	InviteTimeoutMs    int64 `json:"invite_timeout_ms"`
	MatchTimeoutMs     int64 `json:"match_timeout_ms"`
	WaitAttrTimeoutMs  int64 `json:"wait_attr_timeout_ms"`
	ClearRoomTimeoutMs int64 `json:"clear_room_timeout_ms"`
}

func (dtc DelayTimerConfig) InviteTimeout() time.Duration {
	return time.Millisecond * time.Duration(dtc.InviteTimeoutMs)
}

func (dtc DelayTimerConfig) MatchTimeout() time.Duration {
	return time.Millisecond * time.Duration(dtc.MatchTimeoutMs)
}

func (dtc DelayTimerConfig) WaitAttrTimeout() time.Duration {
	return time.Millisecond * time.Duration(dtc.WaitAttrTimeoutMs)
}

func (dtc DelayTimerConfig) ClearRoomTimeout() time.Duration {
	return time.Millisecond * time.Duration(dtc.ClearRoomTimeoutMs)
}
