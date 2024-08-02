package service

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Service interface {
	// CreateGroup creates a new group with the given parameters
	CreateGroup(param *pto.CreateGroup) (entry.Group, error)

	// EnterGroup makes the player join the existed group with the given groupID
	EnterGroup(info *pto.EnterGroup, groupID int64) error

	// ExitGroup makes the player leave the group
	ExitGroup(uid string) error

	// Invite invites the invitee to join the group
	Invite(inviterUID, inviteeUID string) error

	// AcceptInvite accepts the invite and enter the group
	AcceptInvite(inviteeUID string, groupID int64) error

	// RefuseInvite refuses the invite from the inviter
	RefuseInvite(inviteeUID string, groupID int64, refuseMsg string) error

	// StartMatch starts to add the group to matching queue
	StartMatch(captainUID string) error

	// CancelMatch cancels the match and return `entry.GroupStateInvite` state
	CancelMatch(uid string) error

	// UnreadyToMatch makes the player unready to match
	DissolveGroup(uid string) error

	// KickPlayer kicks the kicked player from the group
	KickPlayer(captainUID, kickedUID string) error

	// HandoverCaptain handovers the captain of the group to the target player
	// TODO: is it named ChangeRole better?
	HandoverCaptain(captainUID, targetUID string) error

	// SetNearbyJoinGroup sets whether the group can be joined by nearby players
	SetNearbyJoinGroup(captainUID string, allow bool) error

	// SetRecentJoinGroup sets whether the group can be joined by recent players
	SetRecentJoinGroup(captainUID string, allow bool) error
}
