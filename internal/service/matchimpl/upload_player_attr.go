package matchimpl

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

func (impl *Impl) uploadPlayerAttr(p entry.Player, g entry.Group, attr *pto.UploadPlayerAttr) error {
	if err := p.SetAttr(attr); err != nil {
		return err
	}
	impl.pushService.PushGroupPlayers(g.Base().UIDs(), g.GetPlayerInfos())
	return nil
}
