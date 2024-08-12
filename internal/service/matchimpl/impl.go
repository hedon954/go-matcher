package matchimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/config/mock"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/log"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/internal/repository"
	"github.com/hedon954/go-matcher/internal/service"
	"github.com/hedon954/go-matcher/internal/service/servicemock"
	"github.com/hedon954/go-matcher/pkg/timer"
)

// This file serves as the entry point for the default match service,
// primarily handling validation and preliminary checks.
//
// The detailed business logic is delegated to separate files
// to maintain a clear separation of concerns. For instance,
//  the logic in the `CreateGroup` function, the core logic is extracted into
//  a dedicated `createGroup` function within the `·`create_group.go`·` file.
//
// It is recommended to handle locking and unlocking within this file (most cases),
// while other files should focus mainly on processing the business logic.

// Impl implements a default match service,
// in most cases, you don't need to implement your own match service.
type Impl struct {
	delayTimer timer.Operator[int64]

	DelayConfig config.DelayTimer
	MSConfig    config.MatchStrategy

	playerMgr *repository.PlayerMgr
	groupMgr  *repository.GroupMgr
	teamMgr   *repository.TeamMgr
	roomMgr   *repository.RoomMgr

	groupPlayerLimit int
	nowFunc          func() int64

	// groupChannel used to send a group to match system.
	// TODO: if match system is down, we should stop the server.
	groupChannel chan entry.Group
	roomChannel  chan entry.Room

	pushService        service.Push
	gameServerDispatch service.GameServerDispatch

	result map[int64]*pto.GameResult // TODO: change
}

type Option func(*Impl)

func WithNowFunc(f func() int64) Option {
	return func(impl *Impl) {
		impl.nowFunc = f
	}
}

func WithDelayConfiger(t config.DelayTimer) Option {
	return func(impl *Impl) {
		impl.DelayConfig = t
	}
}

func WithMatchStrategyConfiger(c config.MatchStrategy) Option {
	return func(impl *Impl) {
		impl.MSConfig = c
	}
}

func NewDefault(
	groupPlayerLimit int,
	playerMgr *repository.PlayerMgr, groupMgr *repository.GroupMgr,
	teamMgr *repository.TeamMgr, roomMgr *repository.RoomMgr,
	groupChannel chan entry.Group, roomChannel chan entry.Room,
	delayTimer timer.Operator[int64],
	options ...Option,
) *Impl {
	impl := &Impl{
		playerMgr:          playerMgr,
		groupMgr:           groupMgr,
		teamMgr:            teamMgr,
		roomMgr:            roomMgr,
		groupPlayerLimit:   groupPlayerLimit,
		nowFunc:            time.Now().Unix,
		groupChannel:       groupChannel,
		roomChannel:        roomChannel,
		delayTimer:         delayTimer,                          // TODO: change
		DelayConfig:        new(mock.DelayTimerMock),            // TODO: change
		MSConfig:           new(mock.MatchStrategyMock),         // TODO: change
		pushService:        new(servicemock.PushMock),           // TODO: change
		gameServerDispatch: new(servicemock.GameServerDispatch), // TODO: change
		result:             make(map[int64]*pto.GameResult),     // TODO: change
	}

	for _, opt := range options {
		opt(impl)
	}

	go impl.waitForMatchResult()
	impl.initDelayTimer()
	return impl
}

func (impl *Impl) CreateGroup(ctx context.Context, param *pto.CreateGroup) (entry.Group, error) {
	log.Ctx(ctx).Info().Str("uid", param.UID).Msg("creating group")
	p, err := impl.getPlayer(&param.PlayerInfo)
	if err != nil {
		return nil, err
	}

	p.Base().Lock()
	defer p.Base().Unlock()

	if err := p.Base().CheckOnlineState(entry.PlayerOnlineStateOnline, entry.PlayerOnlineStateInGroup); err != nil {
		return nil, err
	}

	g := impl.groupMgr.Get(p.Base().GroupID)
	if g == nil {
		// create a group
		g, err = impl.createGroup(p)
		if err != nil {
			return nil, err
		}
	} else {
		g.Base().Lock()
		defer g.Base().Unlock()

		// check game mode
		// if game mode is not the same, exits the group and create a new one
		// if game mode is the same, check if the player is the captain of the group
		//  if not, exits the group and create a new one
		//  if yes, return current group
		if g.GetCaptain() != p || g.Base().GameMode != param.GameMode {
			if groupEmpty := g.Base().RemovePlayer(p); groupEmpty {
				impl.groupMgr.Delete(g.ID())
			}
			g, err = impl.createGroup(p)
			if err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}

func (impl *Impl) EnterGroup(ctx context.Context, info *pto.EnterGroup, groupID int64) error {
	g := impl.groupMgr.Get(groupID)
	if g == nil {
		return merr.ErrGroupDissolved
	}

	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	p, err := impl.getPlayer(&info.PlayerInfo)
	if err != nil {
		return err
	}

	p.Base().Lock()
	defer p.Base().Unlock()
	if err := p.Base().CheckOnlineState(entry.PlayerOnlineStateOnline, entry.PlayerOnlineStateInGroup); err != nil {
		return err
	}

	// check source validation
	if err := impl.checkEnterSourceValidation(g, info.Source); err != nil {
		return err
	}

	// player is already in a group
	if p.Base().GroupID != 0 {
		// already in targeted group
		if p.Base().GroupID == groupID {
			// if p can not play together, should exit the origin group
			if err := g.CanPlayTogether(&info.PlayerInfo); err != nil {
				if err := impl.exitGroup(ctx, p, g); err != nil {
					return err
				}
			} else {
				// can play together, refresh the player info and broadcast the group player infos
				p.Base().PlayerInfo = info.PlayerInfo
				impl.pushService.PushGroupInfo(ctx, g.Base().UIDs(), g.GetGroupInfo())
				return nil
			}
		} else {
			// check if player can play together with the group's players
			if err := g.CanPlayTogether(&info.PlayerInfo); err != nil {
				return err
			}

			// not in targeted group, should exit the origin group
			originGroup := impl.groupMgr.Get(p.Base().GroupID)
			if originGroup != nil {
				originGroup.Base().Lock()
				defer originGroup.Base().Unlock()
				if err := impl.exitGroup(ctx, p, originGroup); err != nil {
					return err
				}
			}
		}
	}

	// check if player can play together with the group's players
	if err := g.CanPlayTogether(&info.PlayerInfo); err != nil {
		return err
	}

	// refresh the player info
	p.Base().PlayerInfo = info.PlayerInfo

	// enter the targeted group
	return impl.enterGroup(ctx, p, g)
}

func (impl *Impl) ExitGroup(ctx context.Context, uid string) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	p.Base().Lock()
	defer p.Base().Unlock()
	if err := p.Base().CheckOnlineState(entry.PlayerOnlineStateInGroup); err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()
	if !g.Base().PlayerExists(p.UID()) {
		return nil
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	impl.playerMgr.Delete(p.UID())
	return impl.exitGroup(ctx, p, g)
}

func (impl *Impl) DissolveGroup(ctx context.Context, uid string) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.GetCaptain() != p {
		return merr.ErrOnlyCaptainCanDissolveGroup
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	p.Base().Lock()
	defer p.Base().Unlock()
	return impl.dissolveGroup(ctx, g)
}

func (impl *Impl) KickPlayer(ctx context.Context, captainUID, kickedUID string) error {
	if captainUID == kickedUID {
		return merr.ErrKickSelf
	}

	captain, g, err := impl.getPlayerAndGroup(captainUID)
	if err != nil {
		return err
	}
	kicked := impl.playerMgr.Get(kickedUID)
	if kicked == nil {
		return merr.ErrPlayerNotExists
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.GetCaptain() != captain {
		return merr.ErrOnlyCaptainCanKickPlayer
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	if !g.Base().PlayerExists(kickedUID) {
		return merr.ErrPlayerNotInGroup
	}

	kicked.Base().Lock()
	defer kicked.Base().Unlock()

	impl.kickPlayer(ctx, kicked, g)
	return nil
}

func (impl *Impl) ChangeRole(ctx context.Context, captainUID, targetUID string, role entry.GroupRole) error {
	if captainUID == targetUID {
		return merr.ErrChangeSelfRole
	}

	if err := impl.checkRole(role); err != nil {
		return err
	}

	captain, g, err := impl.getPlayerAndGroup(captainUID)
	if err != nil {
		return err
	}
	target := impl.playerMgr.Get(targetUID)
	if target == nil {
		return merr.ErrPlayerNotExists
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if !g.Base().PlayerExists(targetUID) {
		return merr.ErrPlayerNotInGroup
	}

	if g.GetCaptain() != captain {
		return merr.ErrNotCaptain
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	impl.handoverCaptain(ctx, target, g)
	return nil
}

func (impl *Impl) SetNearbyJoinGroup(_ context.Context, captainUID string, allow bool) error {
	p, g, err := impl.getPlayerAndGroup(captainUID)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.GetCaptain() != p {
		return merr.ErrPermissionDeny
	}

	g.Base().SetAllowNearbyJoin(allow)
	return nil
}

func (impl *Impl) SetRecentJoinGroup(_ context.Context, captainUID string, allow bool) error {
	p, g, err := impl.getPlayerAndGroup(captainUID)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.GetCaptain() != p {
		return merr.ErrPermissionDeny
	}

	g.Base().SetAllowRecentJoin(allow)
	return nil
}

func (impl *Impl) Invite(ctx context.Context, inviterUID, inviteeUID string) error {
	if err := impl.checkInviteeState(inviteeUID); err != nil {
		return err
	}

	inviter, g, err := impl.getPlayerAndGroup(inviterUID)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	if g.IsFull() {
		return merr.ErrGroupFull
	}

	// TODO: how to check if can play together?
	impl.invite(ctx, inviter, inviteeUID, g)
	return nil
}

func (impl *Impl) AcceptInvite(ctx context.Context, inviterUID string, inviteeInfo *pto.PlayerInfo, groupID int64) error {
	g := impl.groupMgr.Get(groupID)
	if g == nil {
		return merr.ErrGroupDissolved
	}

	invitee := impl.playerMgr.Get(inviteeInfo.UID)
	if invitee != nil {
		invitee.Base().Lock()
		defer invitee.Base().Unlock()
		if err := invitee.Base().CheckOnlineState(entry.PlayerOnlineStateOnline,
			entry.PlayerOnlineStateInGroup); err != nil {
			return err
		}
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	// no matter what the result it is, delete the invitation record to the invitee
	defer g.Base().DelInviteRecord(inviteeInfo.UID)

	if !g.Base().PlayerExists(inviterUID) {
		return merr.ErrInvitationExpired
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	if g.IsFull() {
		return merr.ErrGroupFull
	}

	if g.Base().IsInviteExpired(inviteeInfo.UID, impl.nowFunc()) {
		return merr.ErrInvitationExpired
	}

	if err := g.CanPlayTogether(inviteeInfo); err != nil {
		return err
	}

	impl.acceptInvite(ctx, inviterUID, inviteeInfo.UID)
	return nil
}

func (impl *Impl) RefuseInvite(ctx context.Context, inviterUID, inviteeUID string, groupID int64, refuseMsg string) {
	const defaultRefuseMsg = "Sorry, I'm not available at the moment."

	g := impl.groupMgr.Get(groupID)
	if g == nil {
		return
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.Base().GetState() == entry.GroupStateDissolved {
		return
	}
	g.Base().DelInviteRecord(inviteeUID)
	if refuseMsg == "" {
		refuseMsg = defaultRefuseMsg
	}
	impl.pushService.PushRefuseInvite(ctx, inviterUID, inviteeUID, refuseMsg)
}

func (impl *Impl) Ready(ctx context.Context, uid string) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	impl.ready(ctx, p, g)
	return nil
}

func (impl *Impl) Unready(ctx context.Context, uid string) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	impl.unready(ctx, p, g)
	return nil
}

func (impl *Impl) StartMatch(ctx context.Context, captainUID string) error {
	p, g, err := impl.getPlayerAndGroup(captainUID)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	if g.GetCaptain() != p {
		return merr.ErrNotCaptain
	}

	g.Base().MatchStrategy = impl.MSConfig.GetMatchStrategy(g.Base().GameMode)
	if !g.Base().IsMatchStrategySupported() {
		return fmt.Errorf("unsupported match strategy: %v", g.Base().MatchStrategy)
	}

	impl.startMatch(ctx, g)
	return nil
}

func (impl *Impl) CancelMatch(ctx context.Context, uid string) error {
	_, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if err := g.Base().CheckState(entry.GroupStateMatch); err != nil {
		return err
	}

	impl.cancelMatch(ctx, uid, g)
	return nil
}

func (impl *Impl) ExitGame(ctx context.Context, uid string, roomID int64) error {
	panic("implement me")
}

func (impl *Impl) SetVoiceState(ctx context.Context, uid string, state entry.PlayerVoiceState) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	p.Base().Lock()
	defer p.Base().Unlock()

	vs := p.Base().GetVoiceState()
	if vs == state {
		return nil
	}
	p.Base().SetVoiceState(state)

	g.Base().Lock()
	defer g.Base().Unlock()

	impl.pushService.PushVoiceState(ctx, g.Base().UIDs(), &pto.UserVoiceState{UID: uid, State: int(state)})
	return nil
}

func (impl *Impl) UploadPlayerAttr(ctx context.Context, uid string, attr *pto.UploadPlayerAttr) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	p.Base().Lock()
	defer p.Base().Unlock()

	if err := p.Base().CheckOnlineState(entry.PlayerOnlineStateInGroup,
		entry.PlayerOnlineStateInMatch, entry.PlayerOnlineStateInGame); err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	return impl.uploadPlayerAttr(ctx, p, g, attr)
}

func (impl *Impl) HandleMatchResult(r entry.Room) {
	r.Base().Lock()
	defer r.Base().Unlock()
	if err := impl.handleMatchResult(context.Background(), r); err != nil {
		log.Error().
			Any("room", r).
			Err(err).
			Msgf("failed to handle match result: %v", err)
	}
}

func (impl *Impl) HandleGameResult(result *pto.GameResult) error {
	impl.result[result.RoomID] = result
	impl.removeClearRoomTimer(result.RoomID)

	log.Info().
		Int64("room_id", result.RoomID).
		Int("game_mode", int(result.GameMode)).
		Int64("mode_version", result.ModeVersion).
		Int("match_strategy", int(result.MatchStrategy)).
		Any("player_meta_infos", result.PlayerMetaInfo).
		Msg("handle game result")

	// ... do something
	return nil
}

// getPlayerAndGroup returns the player and group of the given uid.
func (impl *Impl) getPlayerAndGroup(uid string) (entry.Player, entry.Group, error) {
	p := impl.playerMgr.Get(uid)
	if p == nil {
		return nil, nil, merr.ErrPlayerNotExists
	}
	g := impl.groupMgr.Get(p.Base().GroupID)
	if g == nil {
		return nil, nil, merr.ErrGroupNotExists
	}
	return p, g, nil
}
