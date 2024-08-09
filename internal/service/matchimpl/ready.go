package matchimpl

import (
	"context"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) ready(ctx context.Context, p entry.Player, g entry.Group) {
	delete(g.Base().UnReadyPlayer, p.UID())
	impl.pushService.PushReady(ctx, g.Base().UIDs(), p.UID())
}
