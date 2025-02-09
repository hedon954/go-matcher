package matchimpl

import (
	"context"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) dissolveGroup(ctx context.Context, g entry.Group) error {
	g.Base().SetState(entry.GroupStateDissolved)

	uids := g.Base().UIDs()
	for _, puid := range g.Base().GetPlayers() {
		p := impl.playerMgr.Get(puid)
		p.Base().SetOnlineStateWithLock(entry.PlayerOnlineStateOnline)
		impl.playerMgr.Delete(p.UID())
	}
	g.Base().ClearPlayers()

	impl.pushService.PushPlayerOnlineState(ctx, uids, entry.PlayerOnlineStateOnline)

	impl.groupMgr.Delete(g.ID())
	impl.pushService.PushGroupDissolve(ctx, uids, g.ID())

	impl.removeInviteTimer(g.ID())
	impl.removeWaitAttrTimer(g.ID())
	impl.removeCancelMatchTimer(g.ID())
	return nil
}
