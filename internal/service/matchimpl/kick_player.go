package matchimpl

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) kickPlayer(ctx context.Context, kicked entry.Player, g entry.Group) {
	logrus.WithFields(logrus.Fields{
		"kicked_uid": kicked.UID(),
		"group_id":   g.ID(),
		"game_mode":  g.Base().GameMode,
	}).Debug("kick player")

	kicked.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
	impl.pushService.PushPlayerOnlineState(ctx, []string{kicked.UID()}, entry.PlayerOnlineStateOnline)
	impl.pushService.PushKick(ctx, kicked.UID(), g.ID())
	impl.playerMgr.Delete(kicked.UID())

	impl.removePlayerFromGroup(ctx, kicked, g)
}

func (impl *Impl) removePlayerFromGroup(ctx context.Context, p entry.Player, g entry.Group) {
	g.Base().RemovePlayer(p)
	impl.pushService.PushGroupInfo(ctx, g.Base().UIDs(), g.GetGroupInfo())
}
