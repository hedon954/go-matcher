package matcher

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Matcher interface {
	CreateGroup(param *pto.CreateGroup) (entry.Group, error)
	EnterGroup(info *pto.PlayerInfo, groupID int64) error
	ExitGroup(uid string) error
	Invite(inviterUID, inviteeUID string) error
	AcceptInvite(inviterUID string, groupID int64) error
	RefuseInvite(inviterUID string, groupID int64, refuseMsg string) error
	CancelMatch(uid string) error
	ReadyToMatch(uid string) error
	DissolveGroup(uid string) error
	KickPlayer(captainUID, kickedUID string) error
	StartMatch(captainUID string) error
	HandoverCaptain(captainUID, targetUID string) error
}
