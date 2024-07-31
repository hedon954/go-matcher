// Package connector provides the functionality to send rpc requests to connector server.
// You can understand the connector servers as a gateway to get requests and push commands to client.
package connector

import (
	"fmt"

	"github.com/hedon954/go-matcher/common"
	"github.com/hedon954/go-matcher/pto"
)

type Manager struct{}

func NewMgr() *Manager {
	return &Manager{}
}

func (m *Manager) PushGroupUsers(uids []string, users []pto.GroupUser) {
	for _, uid := range uids {
		fmt.Print("PushGroupUsers to ", uid)
		for _, user := range users {
			fmt.Println("    ", user)
		}
		fmt.Println()
	}
}

func (m *Manager) PushInviteFriend(param *pto.InviteFriend) {
	fmt.Println("PushInviteFriend to ", param.FriendUid, " ", param)
}

func (m *Manager) PushHandleInvite(inviter string, invitee string, msg int, message string) error {
	fmt.Println("PushHandleInvite to ", inviter, " ", invitee, " ", msg, " ", message)
	return nil
}

func (m *Manager) UpdateOnlineState(uids []string, state common.PlayerOnlineState) {
	for _, uid := range uids {
		fmt.Println("UpdateOnlineState to ", uid, " state: ", state)
	}
}

// 这个 count 是什么意思？
func (m *Manager) UpdateInviteCard(uid string, state common.ChatCardState, count int, src common.InviteCardSrc) {
	fmt.Println("UpdateInviteCard to ", uid, " state: ", state, " count: ", count, " src: ", src)
}

func (m *Manager) GroupDissolved(uids []string, groupID int64) {
	for _, uid := range uids {
		fmt.Println("GroupDissolved to ", uid, " groupID: ", groupID)
	}
}

func (m *Manager) PushGroupState(uids []string, groupID int64, state common.GroupState, name, cancelUID string) {
	for _, uid := range uids {
		fmt.Println("PushGroupState to ", uid, " groupID: ", groupID, " state: ", state, " name: ", name,
			" cancelUID: ", cancelUID)
	}
}

func (m *Manager) CheckInviteFriend(p *pto.CheckInviteFriend) error {
	fmt.Println("CheckInviteFriend to ", p)
	return nil
}

func (m *Manager) PushGroupVoiceState(uids []string, states []*pto.UserVoiceState) {
	for _, uid := range uids {
		fmt.Println("PushGroupVoiceState to ", uid, " ", states)
	}
}
