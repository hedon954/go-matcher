package impl

import (
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) handoverCaptain(captain entry.Player, target entry.Player, g entry.Group) error {
	g.SetCaptain(target)
	impl.connectorClient.PushGroupUsers(g.Base().UIDs(), g.GetPlayerInfos())
	return nil
}
