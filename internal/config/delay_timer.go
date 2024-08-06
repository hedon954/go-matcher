package config

import (
	"time"

	"github.com/hedon954/go-matcher/internal/constant"
)

type DelayTimer interface {
	GetConfig(mode constant.GameMode) DelayTimerConfig
}

// DelayTimerConfig defines the delay timer config.
type DelayTimerConfig struct {
	InviteTimeoutMs   int64 `json:"invite_timeout_ms"`
	MatchTimeoutMs    int64 `json:"match_timeout_ms"`
	WaitAttrTimeoutMs int64 `json:"wait_attr_timeout_ms"`
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
