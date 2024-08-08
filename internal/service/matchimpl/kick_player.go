package matchimpl

import (
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) kickPlayer(kicked entry.Player, g entry.Group) error {
	kicked.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
	impl.pushService.PushPlayerOnlineState([]string{kicked.UID()}, entry.PlayerOnlineStateOnline)
	impl.pushService.PushKick(kicked.UID(), g.ID())
	impl.playerMgr.Delete(kicked.UID())

	impl.removePlayerFromGroup(kicked, g)
	return nil
}

func (impl *Impl) removePlayerFromGroup(p entry.Player, g entry.Group) {
	g.Base().RemovePlayer(p)
	impl.pushService.PushGroupPlayers(g.Base().UIDs(), g.GetPlayerInfos())
}
