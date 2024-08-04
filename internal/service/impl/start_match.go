package impl

import "github.com/hedon954/go-matcher/internal/entry"

func (impl *Impl) startMatch(g entry.Group) {
	base := g.Base()
	uids := base.UIDs()

	// update group state
	base.SetState(entry.GroupStateMatch)
	impl.connectorClient.PushGroupState(uids, g.GroupID(), base.GetState())

	// update players state
	for _, p := range base.GetPlayers() {
		p.Base().Lock()
		p.Base().SetOnlineState(entry.PlayerOnlineStateInMatch)
		p.Base().Unlock()
	}
	impl.connectorClient.UpdateOnlineState(uids, int(entry.PlayerOnlineStateInMatch))

	// push to match channel
	impl.matchChannel <- g
}
