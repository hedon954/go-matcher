package matchimpl

func (impl *Impl) acceptInvite(inviterUID, inviteeUID string) {
	impl.pushService.PushAcceptInvite(inviterUID, inviteeUID)
}
