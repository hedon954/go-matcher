package common

import (
	"testing"

	"matcher/pto"

	"github.com/stretchr/testify/assert"
)

func TestPlayerBase_ShouldWork(t *testing.T) {
	// case: new base
	p := NewPlayerBase(&pto.PlayerInfo{
		Uid: "Uid",
	})
	assert.Equal(t, p, p.Base())
	assert.Equal(t, "Uid", p.UID())
	assert.Equal(t, PlayerOnlineStateOnline, p.GetOnlineState())
	assert.Equal(t, PlayerVoiceStateOff, p.GetVoiceState())

	// case: set base online state should work
	p.SetOnlineState(PlayerOnlineStateGame)
	assert.Equal(t, PlayerOnlineStateGame, p.GetOnlineState())

	// case: set base voice state should work
	p.SetVoiceState(PlayerVoiceStateOff)
	assert.Equal(t, PlayerVoiceStateOff, p.GetVoiceState())
}
