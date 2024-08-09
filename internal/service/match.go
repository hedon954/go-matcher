package service

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Match interface {
	// CreateGroup creates a new group with the given parameters
	CreateGroup(param *pto.CreateGroup) (entry.Group, error)

	// EnterGroup makes the player join the existed group with the given groupID
	EnterGroup(info *pto.EnterGroup, groupID int64) error

	// ExitGroup makes the player leave the group
	ExitGroup(uid string) error

	// DissolveGroup dissolves a group and remove all players of the group
	DissolveGroup(uid string) error

	// Invite invites the invitee to join the group
	Invite(inviterUID, inviteeUID string) error

	// AcceptInvite accepts the invite and enter the group
	AcceptInvite(inviterUID string, inviteeInfo *pto.PlayerInfo, groupID int64) error

	// RefuseInvite refuses the invite from the inviter
	RefuseInvite(inviterUID, inviteeUID string, groupID int64, refuseMsg string)

	// KickPlayer kicks the kicked player from the group
	KickPlayer(captainUID, kickedUID string) error

	// ChangeRole changes the role of the target player
	ChangeRole(captainUID, targetUID string, role entry.GroupRole) error

	// SetNearbyJoinGroup sets whether the group can be joined by nearby players
	SetNearbyJoinGroup(captainUID string, allow bool) error

	// SetRecentJoinGroup sets whether the group can be joined by recent players
	SetRecentJoinGroup(captainUID string, allow bool) error

	// SetVoiceState sets whether the player can speak in the group
	SetVoiceState(uid string, state entry.PlayerVoiceState) error

	// Ready marks the player as ready
	Ready(uid string) error

	// UnReady marks the player as unready
	UnReady(uid string) error

	// StartMatch starts to add the group to matching queue
	StartMatch(captainUID string) error

	// CancelMatch cancels the match and return `entry.GroupStateInvite` state
	CancelMatch(uid string) error

	// UploadPlayerAttr uploads player attributes
	UploadPlayerAttr(uid string, attrs *pto.UploadPlayerAttr) error

	// HandleMatchResult handles the match result
	HandleMatchResult(r entry.Room)

	// HandleGameResult handles the game result
	HandleGameResult(result *pto.GameResult) error
}
