package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase_SetOnlineState(t *testing.T) {
	p := Base{}
	p.SetOnlineState(PlayerOnlineStateGame)
	assert.Equal(t, PlayerOnlineStateGame, p.GetOnlineState())
}
