package player

import (
	"testing"

	"matcher/common"
	"matcher/enum"
	"matcher/pto"

	"github.com/stretchr/testify/assert"
)

const (
	UnknownPlayerUID     = "unknown player uid"
	UnknownGameMode      = 10010
	UnknownMatchStrategy = 10010
)

func Test_Manager_Should_Work(t *testing.T) {
	mgr := NewManager()
	assert.NotNil(t, mgr.players)

	p1UID := "player1"

	// case: create unknown game player should fail
	p, err := mgr.CreatePlayer(&pto.PlayerInfo{
		GameMode: UnknownGameMode,
	})
	assert.NotNil(t, err)
	assert.Nil(t, p)

	// case: create known game but unknown match strategy
	p, err = mgr.CreatePlayer(&pto.PlayerInfo{
		GameMode:      enum.GameModeGoat,
		MatchStrategy: UnknownMatchStrategy,
	})
	assert.NotNil(t, err)
	assert.Nil(t, p)

	// case: create goat game player should work
	p, err = mgr.CreatePlayer(&pto.PlayerInfo{
		GameMode:          enum.GameModeGoat,
		ModeVersion:       1,
		MatchStrategy:     enum.MatchStrategyGlicko2,
		UID:               p1UID,
		UnityNamespacePre: "goat_game",
	})
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, p1UID, p.UID())
	assert.Equal(t, 1, len(mgr.players))
	assert.NotNil(t, p.Inner())
	assert.Equal(t, common.PlayerOnlineStateOnline, p.Inner().GetOnlineState())

	// case: get player
	p1 := mgr.GetPlayer(p1UID)
	assert.Equal(t, p, p1)

	// case: player exists
	assert.True(t, mgr.PlayerExists(p1UID))

	// case: delete player
	mgr.DelPlayer(UnknownPlayerUID)
	assert.Equal(t, 1, len(mgr.players))
	mgr.DelPlayer(p1UID)
	assert.Equal(t, 0, len(mgr.players))
	assert.False(t, mgr.PlayerExists(p1UID))
}
