package matchimpl

import (
	"context"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) exitGroup(ctx context.Context, p entry.Player, g entry.Group) error {
	p.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
	p.Base().GroupID = 0
	empty := g.Base().RemovePlayer(p)
	impl.pushService.PushPlayerOnlineState(ctx, []string{p.UID()}, entry.PlayerOnlineStateOnline)
	if empty {
		return impl.dissolveGroup(ctx, g)
	} else {
		impl.pushService.PushGroupPlayers(ctx, g.Base().UIDs(), g.GetPlayerInfos())
	}
	return nil
}
