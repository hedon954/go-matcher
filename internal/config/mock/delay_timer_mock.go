package mock

import (
	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/constant"
)

const (
	InviteTimeoutMs   = 300000
	MatchTimeoutMs    = 60000
	WaitAttrTimeoutMs = 1
)

type DelayTimerMock struct{}

func (d *DelayTimerMock) GetConfig(_ constant.GameMode) config.DelayTimerConfig {
	return config.DelayTimerConfig{
		InviteTimeoutMs:   InviteTimeoutMs,
		MatchTimeoutMs:    MatchTimeoutMs,
		WaitAttrTimeoutMs: WaitAttrTimeoutMs,
	}
}
