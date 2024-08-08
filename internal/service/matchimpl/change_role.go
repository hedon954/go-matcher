package matchimpl

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) checkRole(role entry.GroupRole) error {
	if role == entry.GroupRoleCaptain {
		return nil
	}
	return fmt.Errorf("unsupported role: %d", role)
}

func (impl *Impl) handoverCaptain(captain entry.Player, target entry.Player, g entry.Group) error {
	g.SetCaptain(target)
	impl.pushService.PushGroupPlayers(g.Base().UIDs(), g.GetPlayerInfos())
	return nil
}
