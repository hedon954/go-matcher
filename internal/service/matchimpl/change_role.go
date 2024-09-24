package matchimpl

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) checkRole(role entry.GroupRole) error {
	if role == entry.GroupRoleCaptain {
		return nil
	}
	return fmt.Errorf("unsupported role: %d", role)
}

func (impl *Impl) handoverCaptain(ctx context.Context, target entry.Player, g entry.Group) {
	logrus.WithFields(logrus.Fields{
		"captain_uid": g.GetCaptain().UID(),
		"target_uid":  target.UID(),
		"group_id":    g.ID(),
		"game_mode":   g.Base().GameMode,
	}).Debug("handover captain")

	g.SetCaptain(target)
	impl.pushService.PushGroupInfo(ctx, g.Base().UIDs(), g.GetGroupInfo())
}
