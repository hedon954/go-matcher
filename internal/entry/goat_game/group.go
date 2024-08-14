package goat_game

import (
	"encoding/json"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

type Group struct {
	*glicko2.GroupBaseGlicko2
}

func CreateGroup(base *entry.GroupBase) entry.Group {
	g := &Group{}
	// ... other common fields

	g.withMatchStrategy(base)
	return g
}

// withMatchStrategy initializes the parameters related to the match strategy.
// We do not initialize the parameters according to the match strategy here,
// because we want to switch the match strategy dynamically without re-initializing the Group object.
func (g *Group) withMatchStrategy(base *entry.GroupBase) {
	g.GroupBaseGlicko2 = glicko2.NewGroup(base)
	// ... other match strategy initialization
}

func (g *Group) Json() string {
	g.Lock()
	defer g.Unlock()
	bs, _ := json.Marshal(g)
	return string(bs)
}
