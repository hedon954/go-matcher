package api

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type EnterGroupReq struct {
	PlayerInfo pto.EnterGroup `json:"player_info" binding:"required"`
	GroupID    int64          `json:"group_id" binding:"required"`
}

type ExitGroupReq struct {
	UID     string `json:"uid" binding:"required"`
	GroupID int64  `json:"group_id" binding:"required"`
}

type KickPlayerReq struct {
	CaptainUID string `json:"captain_uid" binding:"required"`
	KickedUID  string `json:"kicked_uid" binding:"required"`
}

type ChangeRoleReq struct {
	CaptainUID string          `json:"captain_uid" binding:"required"`
	TargetUID  string          `json:"target_uid" binding:"required"`
	Role       entry.GroupRole `json:"role" binding:"required"`
}

type InviteReq struct {
	InviterUID string `json:"inviter_uid" binding:"required"`
	InviteeUID string `json:"invitee_uid" binding:"required"`
}

type AcceptInviteReq struct {
	InviterUID  string          `json:"inviter_uid" binding:"required"`
	InviteeInfo *pto.PlayerInfo `json:"invitee_info" binding:"required"`
	GroupID     int64           `json:"group_id" binding:"required"`
}

type RefuseInviteReq struct {
	InviterUID string `json:"inviter_uid" binding:"required"`
	InviteeUID string `json:"invitee_uid" binding:"required"`
	GroupID    int64  `json:"group_id" binding:"required"`
	RefuseMsg  string `json:"refuse_msg" binding:"optional"`
}

type SetNearbyJoinGroupReq struct {
	CaptainUID string `json:"captain_uid" binding:"required"`
	Allow      bool   `json:"allow"`
}

type SetRecentJoinGroupReq struct {
	CaptainUID string `json:"captain_uid" binding:"required"`
	Allow      bool   `json:"allow"`
}

type SetVoiceStateReq struct {
	UID   string                 `json:"uid" binding:"required"`
	State entry.PlayerVoiceState `json:"state"`
}
