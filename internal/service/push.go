package service

import (
	"context"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

// Push pushes the message to the client.
type Push interface {
	// PushPlayerOnlineState pushes the player's online state to the client.
	PushPlayerOnlineState(ctx context.Context, uids []string, state entry.PlayerOnlineState)

	// PushGroupInfo pushes the group's players infos to the client.
	PushGroupInfo(ctx context.Context, uids []string, users *pto.GroupInfo)

	// PushInviteMsg pushes the invite message to the client.
	PushInviteMsg(ctx context.Context, param *pto.InviteMsg)

	// PushAcceptInvite pushes the accept invite message to the client.
	PushAcceptInvite(ctx context.Context, inviter, invitee string)

	// PushRefuseInvite pushes the refuse invite message to the client.
	PushRefuseInvite(ctx context.Context, inviter, invitee, refuseMsg string)

	// PushGroupDissolve pushes the group dissolve message to the client.
	PushGroupDissolve(ctx context.Context, uids []string, groupID int64)

	// PushGroupState pushes the group state to the client.
	PushGroupState(ctx context.Context, uids []string, groupID int64, state entry.GroupState)

	// PushVoiceState pushes the voice state to the client.
	PushVoiceState(ctx context.Context, uids []string, states *pto.UserVoiceState)

	// PushKick pushes the kick message to the client.
	PushKick(ctx context.Context, uid string, groupID int64)

	// PushMatchInfo pushes the match success info to the client.
	PushMatchInfo(ctx context.Context, uids []string, info *pto.MatchInfo)

	// PushCancelMatch pushes the cancel match message to the client.
	PushCancelMatch(ctx context.Context, uids []string, cancelUID string)

	// PushReady pushes the ready message to the client.
	PushReady(ctx context.Context, uids []string, readyUID string)

	// PushUnReady pushes the unready message to the client.
	PushUnReady(ctx context.Context, uids []string, unreadyUID string)
}
