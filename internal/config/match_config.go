package config

import (
	"time"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type DelayTimerType string

const (
	DelayTimerTypeAsynq  DelayTimerType = "asynq"
	DelayTimerTypeNative DelayTimerType = "native"
)

// MatchConfig defines the global match config.
type MatchConfig struct {
	GroupPlayerLimit int                `yaml:"group_player_limit"`
	MatchIntervalMs  int64              `yaml:"match_interval_ms"`
	Glicko2          *glicko2.QueueArgs `yaml:"glicko2"`
	DelayTimerType   DelayTimerType     `yaml:"delay_timer_type"`
	DelayTimerConfig *DelayTimerConfig  `yaml:"delay_timer_config"`
}

func (c *MatchConfig) GetGlicko2QueueArgs(_ constant.GameMode) *glicko2.QueueArgs {
	return c.Glicko2
}

func (c *MatchConfig) MatchInterval() time.Duration {
	return time.Duration(c.MatchIntervalMs) * time.Millisecond
}
