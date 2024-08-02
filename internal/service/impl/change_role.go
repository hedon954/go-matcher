package impl

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) checkRole(role entry.GroupRole) error {
	switch role {
	case entry.GroupRoleCaptain:
		return nil
	}
	return fmt.Errorf("Unsupported role: %d", role)
}

func (impl *Impl) handoverCaptain(captain entry.Player, target entry.Player, g entry.Group) error {
	g.SetCaptain(target)
	impl.connectorClient.PushGroupUsers(g.Base().UIDs(), g.GetPlayerInfos())
	return nil
}
