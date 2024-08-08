package matchimpl

import (
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) exitGroup(p entry.Player, g entry.Group) error {
	p.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
	p.Base().GroupID = 0
	empty := g.Base().RemovePlayer(p)
	impl.pushService.PushPlayerOnlineState([]string{p.UID()}, entry.PlayerOnlineStateOnline)
	if empty {
		return impl.dissolveGroup(p, g)
	} else {
		impl.pushService.PushGroupPlayers(g.Base().UIDs(), g.GetPlayerInfos())
	}
	return nil
}
