package impl

import (
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) kickPlayer(kicked entry.Player, g entry.Group) error {
	kicked.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
	impl.connectorClient.UpdateOnlineState([]string{kicked.UID()}, int(entry.PlayerOnlineStateOnline))
	impl.connectorClient.PushKick(kicked.UID(), g.GroupID())
	impl.playerMgr.Delete(kicked.UID())

	impl.removePlayerFromGroup(kicked, g)
	return nil
}

func (impl *Impl) removePlayerFromGroup(p entry.Player, g entry.Group) {
	g.Base().RemovePlayer(p)
	impl.connectorClient.PushGroupUsers(g.Base().UIDs(), g.GetPlayerInfos())
}
