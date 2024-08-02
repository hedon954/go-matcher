package entry

import (
	"sync"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
)

// Group represents a group of players.
type Group interface {
	// GroupID returns the unique group id.
	GroupID() int64

	// Base returns the base information of the group.
	// Here we define a concrete struct `GroupBase`
	// to hold the common fields to avoid lots getter and setter method.
	Base() *GroupBase

	// IsFull checks if the group is full.
	IsFull() bool

	// SetCaptain sets the captain of the group.
	SetCaptain(Player)

	// GetCaptain returns the captain in the group.
	GetCaptain() Player

	// CanPlayTogether checks if the player can play with the group's players.
	CanPlayTogether(Player) bool

	// GetPlayerInfos 获取队伍用户信息
	GetPlayerInfos() pto.GroupUser

	// CanStartMatch checks if the group can start to match.
	CanStartMatch() bool
}

const (
	InviteExpireSec = 60 * 5
)

type GroupState int8

const (
	GroupStateInvite    GroupState = 0
	GroupStateMatch     GroupState = 1
	GroupStateGame      GroupState = 2
	GroupStateDissolved GroupState = 3
)

type GroupRole int8

const (
	GroupRoleMember  GroupRole = 0
	GroupRoleCaptain GroupRole = 1
)

// GroupBase holds the common fields of a Group for all kinds of game mode and match strategy.
type GroupBase struct {
	groupID       int64
	GameMode      constant.GameMode
	ModeVersion   int64
	MatchStrategy constant.MatchStrategy

	// Do not do synchronization at this layer,
	// leave it to the caller to handle it uniformly,
	// to avoid deadlocks.
	sync.RWMutex
	state   GroupState
	players []Player

	// MatchID is a unique id to identify each match action.
	MatchID string

	// roles holds the roles of the players in the group.
	roles map[string]GroupRole

	// Settings holds the settings of the group.
	Settings GroupSettings

	// Configs holds the config of the group.
	Configs GroupConfig

	// InviteRecords holds the invite records of the group.
	InviteRecords map[string]int64
}

// GroupSettings defines the settings of a group.
type GroupSettings struct {
	// nearbyJoinAllowed indicates whether nearby players can join the group.
	nearbyJoinAllowed bool

	// necentJoinAllowed indicates whether recent players can join the group.
	necentJoinAllowed bool
}

type GroupConfig struct {
	PlayerLimit     int
	InviteExpireSec int64
}

// NewGroupBase creates a new GroupBase.
func NewGroupBase(
	groupID int64, playerLimit int, mode constant.GameMode, modeVersion int64, strategy constant.MatchStrategy,
) *GroupBase {
	g := &GroupBase{
		groupID:       groupID,
		state:         GroupStateInvite,
		GameMode:      mode,
		ModeVersion:   modeVersion,
		MatchStrategy: strategy,
		players:       make([]Player, 0, playerLimit),
		roles:         make(map[string]GroupRole, playerLimit),
		InviteRecords: make(map[string]int64, playerLimit),
		Configs:       GroupConfig{PlayerLimit: playerLimit, InviteExpireSec: InviteExpireSec},
	}

	return g
}

func (g *GroupBase) Base() *GroupBase {
	return g
}

func (g *GroupBase) GroupID() int64 {
	return g.groupID
}

func (g *GroupBase) IsFull() bool {
	return len(g.players) >= g.Configs.PlayerLimit
}

func (g *GroupBase) GetCaptain() Player {
	uid := ""
	for key, role := range g.roles {
		if role == GroupRoleCaptain {
			uid = key
			break
		}
	}
	for _, p := range g.players {
		if p.UID() == uid {
			return p
		}
	}

	panic("unreachable: group lack of captain")
}

func (g *GroupBase) CanPlayTogether(player Player) bool {
	if g.GameMode != player.Base().GameMode || g.ModeVersion != player.Base().ModeVersion {
		return false
	}
	return true
}

func (g *GroupBase) AddPlayer(p Player) error {
	if g.IsFull() {
		return merr.ErrGroupFull
	}
	if len(g.players) == 0 {
		g.roles[p.UID()] = GroupRoleCaptain
	}
	p.Base().GroupID = g.GroupID()
	p.Base().SetOnlineState(PlayerOnlineStateInGroup)
	for i, player := range g.players {
		if player.UID() == p.UID() {
			g.players[i] = p
			return nil
		}
	}
	g.players = append(g.players, p)
	return nil
}

func (g *GroupBase) GetPlayers() []Player {
	return g.players
}

func (g *GroupBase) UIDs() []string {
	res := make([]string, 0, len(g.players))
	for _, p := range g.players {
		res = append(res, p.UID())
	}
	return res
}

func (g *GroupBase) PlayerLimit() int {
	return g.Configs.PlayerLimit
}

func (g *GroupBase) RemovePlayer(p Player) (empty bool) {
	for index, player := range g.players {
		if player.UID() == p.UID() {
			g.players = append(g.players[:index], g.players[index+1:]...)
			break
		}
	}

	if len(g.players) == 0 {
		g.roles = make(map[string]GroupRole, g.PlayerLimit())
		return true
	} else {
		if g.roles[p.UID()] == GroupRoleCaptain {
			g.SetCaptain(g.players[0])
		}
		delete(g.roles, p.UID())
		return false
	}
}

func (g *GroupBase) ClearPlayers() {
	g.players = make([]Player, 0)
}

func (g *GroupBase) PlayerExists(uid string) bool {
	for _, p := range g.players {
		if p.UID() == uid {
			return true
		}
	}
	return false
}

func (g *GroupBase) SetState(s GroupState) {
	g.state = s
}

func (g *GroupBase) GetState() GroupState {
	return g.state
}

func (g *GroupBase) SetCaptain(p Player) {
	for key, role := range g.roles {
		if role == GroupRoleCaptain {
			g.roles[key] = GroupRoleMember
		}
	}
	g.roles[p.UID()] = GroupRoleCaptain
}

func (g *GroupBase) CheckState(valids ...GroupState) error {
	for _, vs := range valids {
		if g.state == vs {
			return nil
		}
	}

	switch g.state {
	case GroupStateInvite:
		return merr.ErrGroupInInvite
	case GroupStateMatch:
		return merr.ErrGroupInMatch
	case GroupStateGame:
		return merr.ErrGroupInGame
	case GroupStateDissolved:
		return merr.ErrGroupDissolved
	}

	panic("unreachable")
}

func (g *GroupBase) GetPlayerInfos() pto.GroupUser {
	return pto.GroupUser{
		GroupID:  g.groupID,
		Owner:    g.GetCaptain().UID(),
		GameMode: int(g.GameMode),
		// TODO: other info
	}
}

func (g *GroupBase) CanStartMatch() bool {
	return true
}

func (g *GroupBase) SetAllowNearbyJoin(allow bool) {
	g.Settings.nearbyJoinAllowed = allow
}

func (g *GroupBase) SetAllowRecentJoin(allow bool) {
	g.Settings.necentJoinAllowed = allow
}

func (g *GroupBase) AllowNearbyJoin() bool {
	return g.Settings.nearbyJoinAllowed
}

func (g *GroupBase) AllowRecentJoin() bool {
	return g.Settings.necentJoinAllowed
}

func (g *GroupBase) AddInviteRecord(inviteeUID string, nowUnix int64) {
	g.InviteRecords[inviteeUID] = nowUnix + g.Configs.InviteExpireSec
}
