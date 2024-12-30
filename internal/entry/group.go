package entry

import (
	"slices"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/concurrent"
)

// Group represents a group of players.
type Group interface {
	Coder
	Jsoner

	// ID returns the unique group id.
	ID() int64

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
	CanPlayTogether(*pto.PlayerInfo) error

	// GetPlayerInfos returns the player infos of the group.
	// This method usually used for sync group info to client.
	GetGroupInfo() *pto.GroupInfo

	// CanStartMatch checks if the group can start to match.
	// Maybe some game mode need to check if the group is full or not.
	// Maybe some game mode need all players to be ready.
	// If you have some special logics, please override this method.
	CanStartMatch() error

	// GetStartMatchTimeSec returns the start match time of the group.
	GetStartMatchTimeSec() int64

	// SetStartMatchTimeSec sets the start match time of the group.
	SetStartMatchTimeSec(sec int64)

	// Json returns the json string of the group.
	// You may need to lock when marshal it to avoid data race,
	// even if in print log.
	//
	// Note: it should be implemented by the specific game mode entry.
	// TODO: any other great way?
	Json() string
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
	// ReentrantLock is a reentrant lock support multiple locks in the same goroutine.
	// Use it to help avoid deadlock.
	lock *concurrent.ReentrantLock

	// GroupID is the unique id of the group.
	GroupID int64

	// IsAI indicates if the group is an AI group.
	IsAI bool

	// GameMode is the game mode of the group.
	GameMode constant.GameMode

	// ModeVersion is the version of the game mode of the group.
	// Only the same version of the players can be played together.
	ModeVersion int64

	// MatchStrategy is the current match strategy of the group.
	MatchStrategy constant.MatchStrategy

	// SupportMatchStrategies is the supported match strategies of the group.
	SupportMatchStrategies []constant.MatchStrategy

	// State is the current state of the group.
	state GroupState

	// players holds the players in the group.
	players []Player

	// MatchID is a unique id to identify each match action.
	MatchID string

	// StartMatchTimeSec is the start match time of the group.
	StartMatchTimeSec int64

	// roles holds the roles of the players in the group.
	roles map[string]GroupRole

	Configs GroupConfig

	// ReadyPlayer holds the unready players in the group.
	// In the project, we assume that the default does not
	// require all players to prepare to start matching,
	// so the map will be empty by default,
	// and you will need to re-initialize this field in the game mode that requires preparation
	UnReadyPlayer map[string]struct{}

	// inviteRecords holds the invite records of the group.
	// key: uid
	// value: expire time (s)
	inviteRecords map[string]int64

	// Settings holds the settings of the group.
	Settings GroupSettings

	// Configs holds the config of the group.
}

// GroupSettings defines the settings of a group.
type GroupSettings struct {
	// nearbyJoinAllowed indicates whether nearby players can join the group.
	nearbyJoinAllowed bool

	// RecentJoinAllowed indicates whether recent players can join the group.
	recentJoinAllowed bool
}

type GroupConfig struct {
	PlayerLimit     int
	InviteExpireSec int64
}

// NewGroupBase creates a new GroupBase.
func NewGroupBase(
	groupID int64, playerLimit int, playerBase *PlayerBase,
) *GroupBase {
	g := &GroupBase{
		lock:                   new(concurrent.ReentrantLock),
		GroupID:                groupID,
		state:                  GroupStateInvite,
		GameMode:               playerBase.GameMode,
		ModeVersion:            playerBase.ModeVersion,
		players:                make([]Player, 0, playerLimit),
		roles:                  make(map[string]GroupRole, playerLimit),
		inviteRecords:          make(map[string]int64, playerLimit),
		SupportMatchStrategies: make([]constant.MatchStrategy, 0),
		UnReadyPlayer:          make(map[string]struct{}, playerLimit),
		Configs:                GroupConfig{PlayerLimit: playerLimit, InviteExpireSec: InviteExpireSec},
	}

	return g
}

func (g *GroupBase) Base() *GroupBase {
	return g
}

func (g *GroupBase) ID() int64 {
	return g.GroupID
}

// IsMatchStrategySupported checks if the group supports the current match strategy.
func (g *GroupBase) IsMatchStrategySupported() bool {
	return slices.Index(g.SupportMatchStrategies, g.MatchStrategy) >= 0
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

func (g *GroupBase) CanPlayTogether(info *pto.PlayerInfo) error {
	if g.GameMode != info.GameMode {
		return merr.ErrGameModeNotMatch
	}
	if g.ModeVersion < info.ModeVersion {
		return merr.ErrGroupVersionTooLow
	}
	if g.ModeVersion > info.ModeVersion {
		return merr.ErrPlayerVersionTooLow
	}
	return nil
}

func (g *GroupBase) AddPlayer(p Player) error {
	if g.IsFull() {
		return merr.ErrGroupFull
	}
	if len(g.players) == 0 {
		g.roles[p.UID()] = GroupRoleCaptain
	}
	p.Base().GroupID = g.ID()
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

func (g *GroupBase) SetStateWithLock(s GroupState) {
	g.Lock()
	defer g.Unlock()
	g.state = s
}

func (g *GroupBase) GetStateWithLock() GroupState {
	g.Lock()
	defer g.Unlock()
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

func (g *GroupBase) GetGroupInfo() *pto.GroupInfo {
	// TODO: player can change position, set position when crate group or enter group
	// positions := make([]bool, g.PlayerLimit())
	// infos := make([]*pto.GroupPlayerInfo, len(g.players))
	//
	// for _, p := range g.players {
	// 	infos = append(infos, &pto.GroupPlayerInfo{
	// 		UID:   "",
	// 		State: 0,
	// 		Role:  0,
	// 	})
	// }

	return &pto.GroupInfo{
		GroupID:     g.GroupID,
		Captain:     g.GetCaptain().UID(),
		GameMode:    g.GameMode,
		ModeVersion: g.ModeVersion,
		// TODO: other info
	}
}

func (g *GroupBase) CanStartMatch() error {
	return nil
}

func (g *GroupBase) SetAllowNearbyJoin(allow bool) {
	g.Settings.nearbyJoinAllowed = allow
}

func (g *GroupBase) SetAllowRecentJoin(allow bool) {
	g.Settings.recentJoinAllowed = allow
}

func (g *GroupBase) AllowNearbyJoin() bool {
	return g.Settings.nearbyJoinAllowed
}

func (g *GroupBase) AllowRecentJoin() bool {
	return g.Settings.recentJoinAllowed
}

func (g *GroupBase) AddInviteRecord(inviteeUID string, nowUnix int64) {
	g.inviteRecords[inviteeUID] = nowUnix + g.Configs.InviteExpireSec
}

func (g *GroupBase) DelInviteRecord(inviteeUID string) {
	delete(g.inviteRecords, inviteeUID)
}

func (g *GroupBase) GetInviteRecords() map[string]int64 {
	return g.inviteRecords
}
func (g *GroupBase) GetInviteExpireTimeStamp(uid string) int64 {
	return g.inviteRecords[uid]
}
func (g *GroupBase) IsInviteExpired(uid string, nowUnix int64) bool {
	return g.inviteRecords[uid] <= nowUnix
}

func (g *GroupBase) Lock() {
	g.lock.Lock()
}

func (g *GroupBase) Unlock() {
	g.lock.Unlock()
}

func (g *GroupBase) GetStartMatchTimeSec() int64 {
	return g.StartMatchTimeSec
}

func (g *GroupBase) SetStartMatchTimeSec(sec int64) {
	g.StartMatchTimeSec = sec
}
