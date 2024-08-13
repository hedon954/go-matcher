package matchimpl

import (
	"context"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) exitGame(ctx context.Context, p entry.Player, g entry.Group, r entry.Room) error {
	if err := impl.exitGroup(ctx, p, g); err != nil {
		return err
	}
	impl.playerMgr.Delete(p.UID())
	r.Base().AddEscapePlayer(p.UID())
	return nil
}
