package matchimpl

import (
	"context"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) unready(ctx context.Context, p entry.Player, g entry.Group) {
	g.Base().UnReadyPlayer[p.UID()] = struct{}{}
	impl.pushService.PushUnReady(ctx, g.Base().UIDs(), p.UID())
}
