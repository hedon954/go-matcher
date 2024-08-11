package pto

import (
	"github.com/hedon954/go-matcher/internal/constant"
)

// push defines some struct for pushing to client.

// InviteMsg is the invitation message pushed to the client.
type InviteMsg struct {
	InviterUID  string
	InviterName string
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

// GroupInfo is the group players info pushed to the clients to sync group info.
type GroupInfo struct {
	GroupID     int64
	Captain     string
	GameMode    constant.GameMode
	ModeVersion int64

	// Positions indicate whether positions in the room are occupied.
	Positions []bool

	// PlayerInfos holds the player infos, related to the player position.
	// If Positions[i] == false, means PlayerInfos[i] would be nil.
	PlayerInfos []*GroupPlayerInfo
}

type GroupPlayerInfo struct {
	UID         string
	Role        int
	OnlineState int
	VoiceState  int
	Ready       bool
	// ... add more common fields according to your requirement
}

// MatchInfo is the match result pushed to the client.
type MatchInfo struct {
	RoomID          int64
	GameMode        constant.GameMode
	ModeVersion     int64
	MatchStrategy   constant.MatchStrategy
	MatchedTimeUnix int64
	Teams           []MatchTeamInfo
	GameServerInfo  GameServerInfo
}
type MatchTeamInfo struct {
	TeamID  int
	Players []MatchPlayerInfo
}
type MatchPlayerInfo struct {
	UID     string
	GroupID int64
	Attr    Attribute
}

// CancelMatch is the cancel match signal pushed to the client.
type CancelMatch struct {
	CancelUID    string
	CancelReason string
}

// GameServerInfo is the game server info pushed to the client.
type GameServerInfo struct {
	Host     string
	Port     uint16
	Protocol constant.NetProtocol
}
