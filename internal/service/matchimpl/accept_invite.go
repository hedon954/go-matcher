package matchimpl

import "context"

func (impl *Impl) acceptInvite(ctx context.Context, inviterUID, inviteeUID string) {
	impl.pushService.PushAcceptInvite(ctx, inviterUID, inviteeUID)
}
