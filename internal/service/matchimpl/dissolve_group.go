package matchimpl

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) dissolveGroup(ctx context.Context, g entry.Group) error {
	logrus.WithFields(logrus.Fields{
		"group_id":  g.ID(),
		"game_mode": g.Base().GameMode,
	}).Debug("dissolve group")

	g.Base().SetState(entry.GroupStateDissolved)

	uids := g.Base().UIDs()
	for _, p := range g.Base().GetPlayers() {
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
