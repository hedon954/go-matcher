package goat_game

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

type Group struct {
	*glicko2.GroupBaseGlicko2
}

// withMatchStrategy initializes the parameters related to the match strategy.
// We do not initialize the parameters according to the match strategy here,
// because we want to switch the match strategy dynamically without re-initializing the Group object.
func (g *Group) withMatchStrategy(base *entry.GroupBase, playerMgr *entry.PlayerMgr) {
	g.GroupBaseGlicko2 = glicko2.NewGroup(base, playerMgr)
	// ... other match strategy initialization
}

func (g *Group) Json() string {
	return entry.Json(g)
}

func (g *Group) Encode() ([]byte, error) {
	return entry.Encode(g)
}

func (g *Group) Decode(data []byte) error {
	return entry.Decode(data, g)
}
