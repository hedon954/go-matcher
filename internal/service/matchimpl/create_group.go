package matchimpl

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/utils"
)

// getPlayer return player, if not exist, create it.
func (impl *Impl) getPlayer(param *pto.PlayerInfo) (entry.Player, error) {
	p := impl.playerMgr.Get(param.UID)
	if p != nil {
		// TODO: add a refresh player info method
		p.Base().GameMode = param.GameMode
		return p, nil
	}

	return impl.playerMgr.CreatePlayer(param)
}

// createGroup creates group, and adds the player to it,
// current play would be the captain of the group.
func (impl *Impl) createGroup(p entry.Player) (entry.Group, error) {
	logrus.WithFields(logrus.Fields{
		"uid":    p.UID(),
		"player": utils.JsonMarshal(p),
	}).Debug("create group")

	g, err := impl.groupMgr.CreateGroup(impl.Configer.Get().GroupPlayerLimit, p)
	if err != nil {
		return nil, err
	}

	p.Base().GroupID = g.ID()
	p.Base().SetOnlineState(entry.PlayerOnlineStateInGroup)

	impl.addInviteTimer(g.ID(), g.Base().GameMode)

	logrus.WithFields(logrus.Fields{
		"uid": p.UID(),
		"extra": utils.JsonMarshal(gin.H{
			"player": p,
			"group":  g,
		}),
	}).Debug("create group successfully")
	return g, nil
}
