package impl

import (
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) exitGroup(p entry.Player, g entry.Group) error {
	p.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
	p.Base().GroupID = 0
	empty := g.Base().RemovePlayer(p)
	impl.connectorClient.UpdateOnlineState([]string{p.UID()}, int(entry.PlayerOnlineStateOnline))
	if empty {
		return impl.dissolveGroup(g)
	} else {
		impl.connectorClient.PushGroupUsers(g.Base().UIDs(), g.GetPlayerInfos())
	}
	return nil
}
