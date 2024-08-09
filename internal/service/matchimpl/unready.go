package matchimpl

import "github.com/hedon954/go-matcher/internal/entry"

func (impl *Impl) unready(p entry.Player, g entry.Group) {
	g.Base().UnReadyPlayer[p.UID()] = struct{}{}
	impl.pushService.PushUnReady(g.Base().UIDs(), p.UID())
}
