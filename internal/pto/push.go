package pto

import (
	"github.com/hedon954/go-matcher/internal/constant"
)

// push defines some struct for pushing to client.

// InviteMsg is the invitation message pushed to the client.
type InviteMsg struct {
	InviterUID  string
	InviteeUID  string
	Source      EnterGroupSourceType
	GameMode    constant.GameMode
	ModeVersion int64
}

// UserVoiceState is user voice state pushed to the client.
type UserVoiceState struct {
	UID   string
	State int
}

// GroupPlayers is the group players info pushed to the clients to sync group info.
type GroupPlayers struct {
	GroupID     int64
	Captain     string
	GameMode    constant.GameMode
	ModeVersion int64

	// Positions indicate whether positions in the room are occupied.
	Positions []bool

	// Infos holds the player infos, releated to the player position.
	// If Positions[i] == false, means Infos[i] would be nil.
	Infos []*GroupPlayerInfo
}

type GroupPlayerInfo struct {
	UID   string
	State int
	Role  int

	// ... add more common fields according to your requirement
}

// MatchInfo is the match result pushed to the client.
type MatchInfo struct{}

// CancelMatch is the cancel match signal pushed to the client.
type CancelMatch struct {
	CancelUID    string
	CancelReason string
}

// GameServerInfo is the game server info pushed to the client.
type GameServerInfo struct {
	Host     string
	Port     uint16
	Protocal constant.NetProtocal
}
