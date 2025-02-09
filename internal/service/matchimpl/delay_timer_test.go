package matchimpl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hedon954/go-matcher/internal/config/mock"
	"github.com/hedon954/go-matcher/internal/entry"
)

const (
	inviteTimeoutMs    = 5
	matchTimeoutMs     = 1
	waitAttrTimeoutMs  = 1
	clearRoomTimeoutMs = 1
)

func TestImpl_inviteTimeoutHandler_shouldWork(t *testing.T) {
	impl := defaultImpl(1,
		WithMatchStrategyConfiger(new(mock.MatchStrategyMock)),
	)

	g, err := impl.CreateGroup(ctx, newCreateGroupParam(UID))
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetStateWithLock())
	time.Sleep(inviteTimeoutMs + 10*time.Millisecond)
	assert.Equal(t, entry.GroupStateDissolved, g.Base().GetStateWithLock())
}

func TestImpl_matchTimeoutHandler_shouldWork(t *testing.T) {
	impl := defaultImpl(2)

	g, err := impl.CreateGroup(ctx, newCreateGroupParam(UID))
	assert.Nil(t, err)

	err = impl.StartMatch(ctx, g.GetCaptain())
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateMatch, g.Base().GetStateWithLock())
	time.Sleep(matchTimeoutMs + 3*time.Millisecond)
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetStateWithLock())
}

func TestImpl_waitAttrTimeoutHandler_shouldWork(t *testing.T) {
	impl := defaultImpl(1)
	g, err := impl.CreateGroup(ctx, newCreateGroupParam(UID))
	assert.Nil(t, err)
	err = impl.StartMatch(ctx, g.GetCaptain())
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateMatch, g.Base().GetStateWithLock())
	assert.NotNil(t, impl.delayTimer.Get(TimerOpTypeGroupWaitAttr, g.ID()))
	time.Sleep(waitAttrTimeoutMs + 3*time.Millisecond)
	assert.Nil(t, impl.delayTimer.Get(TimerOpTypeGroupWaitAttr, g.ID()))
}
