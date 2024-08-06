package impl

import (
	"log/slog"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/timer"
)

const (
	// TimerOpTypeGroupInvite used to dissolve the group if it not starts game after delay.
	TimerOpTypeGroupInvite timer.OpType = "match:timer_gourp_invite"

	// TimerOpTypeGroupMatchif used to cancel match if the group not matched after delay.
	TimerOpTypeGroupMatch timer.OpType = "match:timer_gourp_match"

	// TimerOpTypeGroupWaitAttr used to wait for players to upload attributes after client clicks `StartMatch`.
	// If all players upload attributes, the group would start to match.
	// If timeout. the group would auto start to match.
	TimerOpTypeGroupWaitAttr timer.OpType = "match:timer_group_wait_attr" // nolint
)

func (impl *Impl) initDelayTimer() {
	impl.delayTimer.Register(TimerOpTypeGroupInvite, impl.inviteTimeoutHandler)
	impl.delayTimer.Register(TimerOpTypeGroupMatch, impl.matchTimeoutHandler)
	impl.delayTimer.Register(TimerOpTypeGroupWaitAttr, impl.waitAttrTimeoutHandler)
}

func (impl *Impl) inviteTimeoutHandler(id int64) {
	g := impl.groupMgr.Get(id)
	if g != nil {
		g.Base().Lock()
		defer g.Base().Unlock()
		if err := impl.dissolveGroup(nil, g); err != nil {
			slog.Error("dissolve group error",
				slog.Int64("group_id", g.Base().GroupID),
				slog.Any("group", g),
				slog.String("err", err.Error()),
			)
		}
	}
}

func (impl *Impl) matchTimeoutHandler(id int64) {
	g := impl.groupMgr.Get(id)
	if g != nil {
		g.Base().Lock()
		defer g.Base().Unlock()
		if g.Base().GetState() == entry.GroupStateMatch {
			impl.cancelMatch("", g)
		}
	}
}

func (impl *Impl) waitAttrTimeoutHandler(id int64) {
	g := impl.groupMgr.Get(id)
	if g != nil {
		g.Base().Lock()
		defer g.Base().Unlock()
		if g.Base().GetState() == entry.GroupStateMatch {
			impl.sendGroupToChannel(g)
		}
	}
}

func (impl *Impl) addInviteTimer(groupID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimerOpTypeGroupInvite, groupID,
		impl.DelayConfig.GetConfig(mode).InviteTimeout())
	if err != nil {
		slog.Error("add invite timer error",
			slog.Int64("groupID", groupID),
			slog.String("err", err.Error()),
		)
	}
}

func (impl *Impl) removeInviteTimer(groupID int64) {
	impl.delayTimer.Remove(TimerOpTypeGroupInvite, groupID)
}

func (impl *Impl) addCancelMatchTimer(groupID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimerOpTypeGroupMatch, groupID,
		impl.DelayConfig.GetConfig(mode).MatchTimeout())
	if err != nil {
		slog.Error("add cancel match timer error",
			slog.Int64("groupID", groupID),
			slog.String("err", err.Error()),
		)
	}
}

func (impl *Impl) removeCancelMatchTimer(groupID int64) {
	impl.delayTimer.Remove(TimerOpTypeGroupMatch, groupID)
}

func (impl *Impl) addWaitAttrTimer(groupID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimerOpTypeGroupWaitAttr, groupID,
		impl.DelayConfig.GetConfig(mode).WaitAttrTimeout())
	if err != nil {
		slog.Error("add wait attr timer error",
			slog.Int64("groupID", groupID),
			slog.String("err", err.Error()),
		)
	}
}

func (impl *Impl) removeWaitAttrTimer(groupID int64) {
	impl.delayTimer.Remove(TimerOpTypeGroupWaitAttr, groupID)
}
