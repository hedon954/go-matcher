package matchimpl

import (
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) dissolveGroup(player entry.Player, g entry.Group) error {
	g.Base().SetState(entry.GroupStateDissolved)

	uids := g.Base().UIDs()
	for _, p := range g.Base().GetPlayers() {
		// to avoid deadlock
		if player != nil && p.UID() == player.UID() {
			p.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
		} else {
			p.Base().SetOnlineStateWithLock(entry.PlayerOnlineStateOnline)
		}
		impl.playerMgr.Delete(p.UID())
	}
	g.Base().ClearPlayers()

	impl.pushService.PushPlayerOnlineState(uids, entry.PlayerOnlineStateOnline)

	impl.groupMgr.Delete(g.ID())
	impl.pushService.PushGroupDissolve(uids, g.ID())

	impl.removeInviteTimer(g.ID())
	impl.removeWaitAttrTimer(g.ID())
	impl.removeCancelMatchTimer(g.ID())
	return nil
}
