package player

import (
	"testing"

	"matcher/enum"
	"matcher/pto"

	"github.com/stretchr/testify/assert"
)

func Test_Manager_Should_Work(t *testing.T) {
	mgr := NewManager()
	assert.NotNil(t, mgr.players)

	p1UID := "player1"

	// case1: create unknown game player should fail
	p, err := mgr.CreatePlayer(&pto.PlayerInfo{
		GameMode: 10010,
	})
	assert.NotNil(t, err)
	assert.Nil(t, p)

	// case2: create goat game player should work
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

	// case3: get player
	p1 := mgr.GetPlayer(p1UID)
	assert.Equal(t, p, p1)

	// case4: player exists
	assert.True(t, mgr.PlayerExists(p1UID))

	// case5: delete player
	mgr.DelPlayer("unknown player uid")
	assert.Equal(t, 1, len(mgr.players))
	mgr.DelPlayer(p1UID)
	assert.Equal(t, 0, len(mgr.players))
	assert.False(t, mgr.PlayerExists(p1UID))
}
