package matchimpl

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) exitGroup(ctx context.Context, p entry.Player, g entry.Group) error {
	logrus.WithFields(logrus.Fields{
		"uid":       p.UID(),
		"group_id":  g.ID(),
		"game_mode": g.Base().GameMode,
	}).Info("exit group")

	p.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
	p.Base().GroupID = 0
	empty := g.Base().RemovePlayer(p)
	impl.pushService.PushPlayerOnlineState(ctx, []string{p.UID()}, entry.PlayerOnlineStateOnline)
	if empty {
		return impl.dissolveGroup(ctx, g)
	} else {
		impl.pushService.PushGroupInfo(ctx, g.Base().UIDs(), g.GetGroupInfo())
	}
	return nil
}
