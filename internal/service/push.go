package service

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

// Push pushes the message to the client.
type Push interface {
	// PushPlayerOnlineState pushes the player's online state to the client.
	PushPlayerOnlineState(uids []string, state entry.PlayerOnlineState)

	// PushGroupPlayers pushes the group's players infos to the client.
	PushGroupPlayers(uids []string, users *pto.GroupPlayers)

	// PushInviteMsg pushes the invite message to the client.
	PushInviteMsg(param *pto.InviteMsg)

	// PushAcceptInvite pushes the accept invite message to the client.
	PushAcceptInvite(inviter, invitee string)

	// PushRefuseInvite pushes the refuse invite message to the client.
	PushRefuseInvite(inviter, invitee, refuseMsg string)

	// PushGroupDissolve pushes the group dissolve message to the client.
	PushGroupDissolve(uids []string, groupID int64)

	// PushGroupState pushes the group state to the client.
	PushGroupState(uids []string, groupID int64, state entry.GroupState)

	// PushVoiceState pushes the voice state to the client.
	PushVoiceState(uids []string, states *pto.UserVoiceState)

	// PushKick pushes the kick message to the client.
	PushKick(uid string, groupID int64)

	// PushMatchInfo pushes the match success info to the client.
	PushMatchInfo(uids []string, info *pto.MatchInfo)

	// PushCancelMatch pushes the cancel match message to the client.
	PushCancelMatch(uids []string, cancelUID string)

	// PushReady pushes the ready message to the client.
	PushReady(uids []string, readyUID string)

	// PushUnReady pushes the unready message to the client.
	PushUnReady(uids []string, readyUID string)
}
