package impl

import (
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) enterGroup(p entry.Player, g entry.Group) error {
	if err := g.Base().AddPlayer(p); err != nil {
		return err
	}
	p.Base().GroupID = g.GroupID()
	p.Base().SetOnlineState(entry.PlayerOnlineStateInGroup)
	impl.connectorClient.UpdateOnlineState([]string{p.UID()}, int(entry.PlayerOnlineStateInGroup))
	impl.connectorClient.PushGroupUsers(g.Base().UIDs(), g.GetPlayerInfos())
	return nil
}
