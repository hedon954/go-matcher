package impl

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

// getPlayer return player, if not exist, create it.
func (impl *Impl) getPlayer(param *pto.PlayerInfo) (entry.Player, error) {
	p := impl.playerMgr.Get(param.UID)
	if p != nil {
		// TODO: add a refresh player info method
		p.Base().GameMode = param.GameMode
		p.Base().MatchStrategy = param.MatchStrategy
		return p, nil
	}

	return impl.playerMgr.CreatePlayer(param)
}

// createGroup creates group, and adds the player to it,
// current play would be the captain of the group.
func (impl *Impl) createGroup(param *pto.CreateGroup, p entry.Player) (entry.Group, error) {
	g, err := impl.groupMgr.CreateGroup(impl.playerLimit, param.GameMode, param.ModeVersion, param.MatchStrategy)
	if err != nil {
		return nil, err
	}

	if err := g.Base().AddPlayer(p); err != nil {
		panic("add player to group failed when create group")
	}

	p.Base().GroupID = g.ID()
	p.Base().SetOnlineState(entry.PlayerOnlineStateInGroup)
	return g, nil
}
