package matchimpl

import (
	"testing"
	"time"

	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/config/mock"
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/stretchr/testify/assert"
)

type delayConfiger struct{}

const (
	inviteTimeoutMs   = 5
	matchTimeoutMs    = 1
	waitAttrTimeoutMs = 1
)

func (d *delayConfiger) GetConfig(mode constant.GameMode) config.DelayTimerConfig {
	return config.DelayTimerConfig{
		InviteTimeoutMs:   inviteTimeoutMs,
		MatchTimeoutMs:    matchTimeoutMs,
		WaitAttrTimeoutMs: waitAttrTimeoutMs,
	}
}

func TestImpl_inviteTimeoutHandler_shouldwork(t *testing.T) {
	impl := defaultImpl(1,
		WithDelayConfiger(new(delayConfiger)),
		WithMatchStrategyConfiger(new(mock.MatchStrategyMock)),
	)

	g, err := impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetStateWithLock())
	time.Sleep(inviteTimeoutMs + 10*time.Millisecond)
	assert.Equal(t, entry.GroupStateDissolved, g.Base().GetStateWithLock())
}

func TestImpl_matchTimeoutHandler_shouldwork(t *testing.T) {
	impl := defaultImpl(2, WithDelayConfiger(new(delayConfiger)))

	g, err := impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)

	err = impl.StartMatch(g.GetCaptain().UID())
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateMatch, g.Base().GetStateWithLock())
	time.Sleep(matchTimeoutMs + 3*time.Millisecond)
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetStateWithLock())
}

func TestImpl_waitAttrTimeoutHandler_shouldwork(t *testing.T) {
	impl := defaultImpl(1, WithDelayConfiger(new(delayConfiger)))
	g, err := impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)
	err = impl.StartMatch(g.GetCaptain().UID())
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateMatch, g.Base().GetStateWithLock())
	assert.NotNil(t, impl.delayTimer.Get(TimerOpTypeGroupWaitAttr, g.ID()))
	time.Sleep(waitAttrTimeoutMs + 3*time.Millisecond)
	assert.Nil(t, impl.delayTimer.Get(TimerOpTypeGroupWaitAttr, g.ID()))
}
