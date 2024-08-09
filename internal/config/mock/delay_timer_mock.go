package mock

import (
	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/constant"
)

const (
	InviteTimeoutMs    = 300 * 1000
	MatchTimeoutMs     = 60 * 1000
	WaitAttrTimeoutMs  = 1
	ClearRoomTimeoutMs = 30 * 60 * 1000
)

type DelayTimerMock struct{}

func (d *DelayTimerMock) GetConfig(_ constant.GameMode) config.DelayTimerConfig {
	return config.DelayTimerConfig{
		InviteTimeoutMs:    InviteTimeoutMs,
		MatchTimeoutMs:     MatchTimeoutMs,
		WaitAttrTimeoutMs:  WaitAttrTimeoutMs,
		ClearRoomTimeoutMs: ClearRoomTimeoutMs,
	}
}
