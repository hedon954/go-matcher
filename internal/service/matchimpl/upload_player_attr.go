package matchimpl

import (
	"context"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

func (impl *Impl) uploadPlayerAttr(ctx context.Context, p entry.Player, g entry.Group, attr *pto.UploadPlayerAttr) error {
	if err := p.SetAttr(attr); err != nil {
		return err
	}
	impl.pushService.PushGroupPlayers(ctx, g.Base().UIDs(), g.GetPlayerInfos())
	return nil
}
