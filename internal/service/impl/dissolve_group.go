package impl

import (
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) dissolveGroup(g entry.Group) error {
	g.Base().SetState(entry.GroupStateDissolved)

	uids := g.Base().UIDs()
	for _, p := range g.Base().GetPlayers() {
		p.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
		impl.playerMgr.Delete(p.UID())
	}
	g.Base().ClearPlayers()

	impl.connectorClient.UpdateOnlineState(uids, int(entry.PlayerOnlineStateOnline))

	impl.groupMgr.Delete(g.ID())
	impl.connectorClient.GroupDissolved(uids, g.ID())
	return nil
}
