// Package connector provides the functionality to send rpc requests to connector server.
// You can understand the connector servers as a gateway to get requests and push commands to client.
package connector

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (m *Client) PushGroupUsers(uids []string, users pto.GroupUser) {
	for _, uid := range uids {
		fmt.Print("PushGroupUsers to ", uid, ": ", users)
	}
}

func (m *Client) PushInviteMsg(param *pto.InviteMsg) {
	fmt.Println("PushInviteMsg: ", param)
}

func (m *Client) PushAcceptInvite(inviter, invitee string) {
	fmt.Println("PushAcceptInvite: ", inviter, invitee)
}

func (m *Client) PushRefuseInvite(inviter, invitee, refuseMsg string) {
	fmt.Println("PushRefuseInvite: ", inviter, invitee, refuseMsg)
}

func (m *Client) UpdateOnlineState(uids []string, state int) {
	for _, uid := range uids {
		fmt.Println("UpdateOnlineState to ", uid, " state: ", state)
	}
}

// // 这个 count 是什么意思？
// func (m *Client) UpdateInviteCard(uid string, state entry.ChatCardState, count int, src entry.InviteCardSrc) {
// 	fmt.Println("UpdateInviteCard to ", uid, " state: ", state, " count: ", count, " src: ", src)
// }

func (m *Client) GroupDissolved(uids []string, groupID int64) {
	for _, uid := range uids {
		fmt.Println("GroupDissolved to ", uid, " groupID: ", groupID)
	}
}

func (m *Client) PushGroupState(uids []string, groupID int64, state entry.GroupState, name, cancelUID string) {
	for _, uid := range uids {
		fmt.Println("PushGroupState to ", uid, " groupID: ", groupID, " state: ", state, " name: ", name,
			" cancelUID: ", cancelUID)
	}
}

func (m *Client) CheckInviteFriend(p *pto.CheckInviteFriend) error {
	fmt.Println("CheckInviteFriend to ", p)
	return nil
}

func (m *Client) PushGroupVoiceState(uids []string, states []*pto.UserVoiceState) {
	for _, uid := range uids {
		fmt.Println("PushGroupVoiceState to ", uid, " ", states)
	}
}

func (m *Client) PushKick(uid string, groupID int64) {
	fmt.Println("PushKick to ", uid, " groupID: ", groupID)
}
