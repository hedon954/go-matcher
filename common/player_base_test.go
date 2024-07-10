package common

import (
	"testing"

	"matcher/pto"

	"github.com/stretchr/testify/assert"
)

func TestBase_ShouldWork(t *testing.T) {
	// case: new base
	p := NewBase(&pto.PlayerInfo{
		UID: "uid",
	})
	assert.Equal(t, p, p.Inner())
	assert.Equal(t, "uid", p.UID())
	assert.Equal(t, PlayerOnlineStateOnline, p.GetOnlineState())
	assert.Equal(t, PlayerVoiceStateOff, p.GetVoiceState())

	// case: set base online state should work
	p.SetOnlineState(PlayerOnlineStateGame)
	assert.Equal(t, PlayerOnlineStateGame, p.GetOnlineState())

	// case: set base voice state should work
	p.SetVoiceState(PlayerVoiceStateOff)
	assert.Equal(t, PlayerVoiceStateOff, p.GetVoiceState())
}
