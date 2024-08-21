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

type MatchConfigerMock struct {
	mc *config.MatchConfig
}

func NewMatchConfigerMock(c *config.MatchConfig) *MatchConfigerMock {
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
			RoomTeamLimit:   2, //nolint:mnd
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
	return &MatchConfigerMock{mc: c}
}

func (cm *MatchConfigerMock) Get() *config.MatchConfig {
	return cm.mc
}
