package matchimpl

import (
	"context"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/timer"
	"github.com/rs/zerolog/log"
)

const (
	// TimerOpTypeGroupInvite used to dissolve the group if it not starts game after delay.
	TimerOpTypeGroupInvite timer.OpType = "match:timer_group_invite"

	// TimerOpTypeGroupMatch if used to cancel match if the group not matched after delay.
	TimerOpTypeGroupMatch timer.OpType = "match:timer_group_match"

	// TimerOpTypeGroupWaitAttr used to wait for players to upload attributes after client clicks `StartMatch`.
	// If all players upload attributes, the group would start to match.
	// If timeout. the group would auto start to match.
	TimerOpTypeGroupWaitAttr timer.OpType = "match:timer_group_wait_attr" // nolint

	// TimeOpTypeClearRoom used to clear room in some unexpected cases like client do not settle game.
	// We use this optype to force clear the room info.
	TimeOpTypeClearRoom timer.OpType = "match:timer_clear_room"
)

func (impl *Impl) initDelayTimer() {
	impl.delayTimer.Register(TimerOpTypeGroupInvite, impl.inviteTimeoutHandler)
	impl.delayTimer.Register(TimerOpTypeGroupMatch, impl.matchTimeoutHandler)
	impl.delayTimer.Register(TimerOpTypeGroupWaitAttr, impl.waitAttrTimeoutHandler)
	impl.delayTimer.Register(TimeOpTypeClearRoom, impl.clearRoomTimeoutHandler)
}

func (impl *Impl) inviteTimeoutHandler(groupID int64) {
	g := impl.groupMgr.Get(groupID)
	if g != nil {
		g.Base().Lock()
		defer g.Base().Unlock()
		if err := impl.dissolveGroup(context.Background(), g); err != nil {
			log.Error().
				Int64("group_id", g.Base().GroupID).
				Any("group", g).
				Err(err).
				Msg("dissolve group error")
		}
	}
}

func (impl *Impl) matchTimeoutHandler(groupID int64) {
	g := impl.groupMgr.Get(groupID)
	if g != nil {
		g.Base().Lock()
		defer g.Base().Unlock()
		if g.Base().GetState() == entry.GroupStateMatch {
			impl.cancelMatch(context.Background(), "", g)
		}
	}
}

func (impl *Impl) waitAttrTimeoutHandler(groupID int64) {
	g := impl.groupMgr.Get(groupID)
	if g != nil {
		g.Base().Lock()
		defer g.Base().Unlock()
		if g.Base().GetState() == entry.GroupStateMatch {
			impl.sendGroupToChannel(g)
		}
	}
}

func (impl *Impl) clearRoomTimeoutHandler(roomID int64) {
	r := impl.roomMgr.Get(roomID)
	if r != nil {
		impl.roomMgr.Delete(roomID)
		log.Warn().
			Int64("room_id", roomID).
			Any("room_info", r).
			Msg("clear room timeout")
	}
}

func (impl *Impl) addInviteTimer(groupID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimerOpTypeGroupInvite, groupID,
		impl.DelayConfig.GetConfig(mode).InviteTimeout())
	if err != nil {
		log.Error().
			Int64("group_id", groupID).
			Int("mode", int(mode)).
			Err(err).
			Msg("add invite timer error")
	}
}

func (impl *Impl) removeInviteTimer(groupID int64) {
	_ = impl.delayTimer.Remove(TimerOpTypeGroupInvite, groupID)
}

func (impl *Impl) addCancelMatchTimer(groupID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimerOpTypeGroupMatch, groupID,
		impl.DelayConfig.GetConfig(mode).MatchTimeout())
	if err != nil {
		log.Error().
			Int64("group_id", groupID).
			Int("mode", int(mode)).
			Err(err).
			Msg("add cancel match timer error")
	}
}

func (impl *Impl) removeCancelMatchTimer(groupID int64) {
	_ = impl.delayTimer.Remove(TimerOpTypeGroupMatch, groupID)
}

func (impl *Impl) addWaitAttrTimer(groupID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimerOpTypeGroupWaitAttr, groupID,
		impl.DelayConfig.GetConfig(mode).WaitAttrTimeout())
	if err != nil {
		log.Error().
			Int64("group_id", groupID).
			Int("mode", int(mode)).
			Err(err).
			Msg("add wait attr timer error")
	}
}

func (impl *Impl) removeWaitAttrTimer(groupID int64) {
	_ = impl.delayTimer.Remove(TimerOpTypeGroupWaitAttr, groupID)
}

func (impl *Impl) addClearRoomTimer(roomID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimeOpTypeClearRoom, roomID,
		impl.DelayConfig.GetConfig(mode).ClearRoomTimeout())
	if err != nil {
		log.Error().
			Int64("room_id", roomID).
			Int("mode", int(mode)).
			Err(err).
			Msg("add clear room timer error")
	}
}

func (impl *Impl) removeClearRoomTimer(roomID int64) {
	_ = impl.delayTimer.Remove(TimeOpTypeClearRoom, roomID)
}
