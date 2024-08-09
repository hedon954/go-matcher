package matchimpl

import "github.com/hedon954/go-matcher/internal/entry"

func (impl *Impl) ready(p entry.Player, g entry.Group) {
	delete(g.Base().UnReadyPlayer, p.UID())
	impl.pushService.PushReady(g.Base().UIDs(), p.UID())
}
