package matchimpl

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
)

func (impl *Impl) checkEnterSourceValidation(g entry.Group, source pto.EnterGroupSourceType) error {
	switch source {
	case pto.EnterGroupSourceTypeNearby:
		if !g.Base().AllowNearbyJoin() {
			return merr.ErrGroupDenyNearbyJoin
		}
	case pto.EnterGroupSourceTypeRecent:
		if !g.Base().AllowRecentJoin() {
			return merr.ErrGroupDenyRecentJoin
		}
		// TODO: check other
	}

	return nil
}

func (impl *Impl) enterGroup(ctx context.Context, p entry.Player, g entry.Group) error {
	logrus.WithFields(logrus.Fields{
		"uid":       p.UID(),
		"group_id":  g.ID(),
		"game_mode": g.Base().GameMode,
	}).Debug("enter group")

	if err := g.Base().AddPlayer(p); err != nil {
		return err
	}
	impl.pushService.PushPlayerOnlineState(ctx, []string{p.UID()}, entry.PlayerOnlineStateInGroup)
	impl.pushService.PushGroupInfo(ctx, g.Base().UIDs(), g.GetGroupInfo())
	return nil
}
