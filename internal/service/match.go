package service

import (
	"context"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/matcher/common"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Match interface {
	// CreateGroup creates a new group with the given parameters
	CreateGroup(ctx context.Context, param *pto.CreateGroup) (entry.Group, error)

	// EnterGroup makes the player join the existed group with the given groupID
	EnterGroup(ctx context.Context, info *pto.EnterGroup, groupID int64) error

	// ExitGroup makes the player leave the group
	ExitGroup(ctx context.Context, uid string) error

	// DissolveGroup dissolves a group and remove all players of the group
	DissolveGroup(ctx context.Context, uid string) error

	// Invite invites the invitee to join the group
	Invite(ctx context.Context, inviterUID, inviteeUID string) error

	// AcceptInvite accepts the invite and enter the group
	AcceptInvite(ctx context.Context, inviterUID string, inviteeInfo *pto.PlayerInfo, groupID int64) error

	// RefuseInvite refuses the invite from the inviter
	RefuseInvite(ctx context.Context, inviterUID, inviteeUID string, groupID int64, refuseMsg string)

	// KickPlayer kicks the kicked player from the group
	KickPlayer(ctx context.Context, captainUID, kickedUID string) error

	// ChangeRole changes the role of the target player
	ChangeRole(ctx context.Context, captainUID, targetUID string, role entry.GroupRole) error

	// SetNearbyJoinGroup sets whether the group can be joined by nearby players
	SetNearbyJoinGroup(ctx context.Context, captainUID string, allow bool) error

	// SetRecentJoinGroup sets whether the group can be joined by recent players
	SetRecentJoinGroup(ctx context.Context, captainUID string, allow bool) error

	// SetVoiceState sets whether the player can speak in the group
	SetVoiceState(ctx context.Context, uid string, state entry.PlayerVoiceState) error

	// Ready marks the player as ready
	Ready(ctx context.Context, uid string) error

	// UnReady marks the player as unready
	Unready(ctx context.Context, uid string) error

	// StartMatch starts to add the group to matching queue
	StartMatch(ctx context.Context, captainUID string) error

	// CancelMatch cancels the match and return `entry.GroupStateInvite` state
	CancelMatch(ctx context.Context, uid string) error

	// ExitGame exits the game (escape)
	ExitGame(ctx context.Context, uid string, roomID int64) error

	// UploadPlayerAttr uploads player attributes
	UploadPlayerAttr(ctx context.Context, uid string, attrs *pto.UploadPlayerAttr) error

	// HandleMatchResult handles the match result
	HandleMatchResult(r common.Result)

	// HandleGameResult handles the game result
	HandleGameResult(result *pto.GameResult) error
}
