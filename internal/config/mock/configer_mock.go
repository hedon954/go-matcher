package mock

import (
	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

const (
	InviteTimeoutMs    = 300 * 1000
	MatchTimeoutMs     = 60 * 1000
	WaitAttrTimeoutMs  = 1
	ClearRoomTimeoutMs = 30 * 60 * 1000

	MatchTimeoutSec = MatchTimeoutMs / 1000
)

type ConfigerMock struct {
	config *config.Config
}

func NewConfigerMock(c *config.Config) *ConfigerMock {
	if c.GroupPlayerLimit == 0 {
		c.GroupPlayerLimit = 0
	}
	if c.MatchIntervalMs == 0 {
		c.MatchIntervalMs = 1
	}
	if c.Glicko2 == nil {
		c.Glicko2 = &glicko2.QueueArgs{
			MatchTimeoutSec: MatchTimeoutSec,
			TeamPlayerLimit: c.GroupPlayerLimit,
			RoomTeamLimit:   1,
		}
	}
	if c.DelayTimerType == "" {
		c.DelayTimerType = config.DelayTimerTypeNative
	}
	if c.DelayTimerConfig == nil {
		c.DelayTimerConfig = &config.DelayTimerConfig{
			InviteTimeoutMs:    InviteTimeoutMs,
			MatchTimeoutMs:     MatchTimeoutMs,
			WaitAttrTimeoutMs:  WaitAttrTimeoutMs,
			ClearRoomTimeoutMs: ClearRoomTimeoutMs,
		}
	}
	return &ConfigerMock{config: c}
}

func (cm *ConfigerMock) Get() *config.Config {
	return cm.config
}
