package impl

func (impl *Impl) acceptInvite(inviterUID, inviteeUID string) {
	impl.connectorClient.PushAcceptInvite(inviterUID, inviteeUID)
}
