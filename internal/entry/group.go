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
	SetCaptain(string)

	// GetCaptain returns the captain in the group.
	GetCaptain() string

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
	// TODO: any other greater way?
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
	// ReentrantLock is a reentrant L support multiple locks in the same goroutine.
	// Use it to help avoid deadlock.
	L *concurrent.ReentrantLock

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

	// State is the current State of the group.
	State GroupState

	// Players holds the Players'ids in the group.
	Players []string

	// MatchID is a unique id to identify each match action.
	MatchID string

	// StartMatchTimeSec is the start match time of the group.
	StartMatchTimeSec int64

	// Roles holds the Roles of the players in the group.
	Roles map[string]GroupRole

	Configs GroupConfig

	// ReadyPlayer holds the unready players in the group.
	// In the project, we assume that the default does not
	// require all players to prepare to start matching,
	// so the map will be empty by default,
	// and you will need to re-initialize this field in the game mode that requires preparation
	UnReadyPlayer map[string]struct{}

	// InviteRecords holds the invite records of the group.
	// key: uid
	// value: expire time (s)
	InviteRecords map[string]int64

	// Settings holds the settings of the group.
	Settings GroupSettings

	// Configs holds the config of the group.
}

// GroupSettings defines the settings of a group.
type GroupSettings struct {
	// NearbyJoinAllowed indicates whether nearby players can join the group.
	NearbyJoinAllowed bool

	// RecentJoinAllowed indicates whether recent players can join the group.
	RecentJoinAllowed bool
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
		L:                      new(concurrent.ReentrantLock),
		GroupID:                groupID,
		State:                  GroupStateInvite,
		GameMode:               playerBase.GameMode,
		ModeVersion:            playerBase.ModeVersion,
		Players:                make([]string, 0, playerLimit),
		Roles:                  make(map[string]GroupRole, playerLimit),
		InviteRecords:          make(map[string]int64, playerLimit),
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
	return len(g.Players) >= g.Configs.PlayerLimit
}

func (g *GroupBase) GetCaptain() string {
	uid := ""
	for key, role := range g.Roles {
		if role == GroupRoleCaptain {
			uid = key
			break
		}
	}
	for _, puid := range g.Players {
		if puid == uid {
			return puid
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
	if len(g.Players) == 0 {
		g.Roles[p.UID()] = GroupRoleCaptain
	}
	p.Base().GroupID = g.ID()
	p.Base().SetOnlineState(PlayerOnlineStateInGroup)
	for i, puid := range g.Players {
		if puid == p.UID() {
			g.Players[i] = p.UID()
			return nil
		}
	}
	g.Players = append(g.Players, p.UID())
	return nil
}

func (g *GroupBase) GetPlayers() []string {
	return g.Players
}

func (g *GroupBase) UIDs() []string {
	return g.Players
}

func (g *GroupBase) PlayerLimit() int {
	return g.Configs.PlayerLimit
}

func (g *GroupBase) RemovePlayer(p Player) (empty bool) {
	for index, puid := range g.Players {
		if puid == p.UID() {
			g.Players = append(g.Players[:index], g.Players[index+1:]...)
			break
		}
	}

	if len(g.Players) == 0 {
		g.Roles = make(map[string]GroupRole, g.PlayerLimit())
		return true
	} else {
		if g.Roles[p.UID()] == GroupRoleCaptain {
			g.SetCaptain(g.Players[0])
		}
		delete(g.Roles, p.UID())
		return false
	}
}

func (g *GroupBase) ClearPlayers() {
	g.Players = make([]string, 0)
}

func (g *GroupBase) PlayerExists(uid string) bool {
	for _, puid := range g.Players {
		if puid == uid {
			return true
		}
	}
	return false
}

func (g *GroupBase) SetState(s GroupState) {
	g.State = s
}

func (g *GroupBase) GetState() GroupState {
	return g.State
}

func (g *GroupBase) SetStateWithLock(s GroupState) {
	g.Lock()
	defer g.Unlock()
	g.State = s
}

func (g *GroupBase) GetStateWithLock() GroupState {
	g.Lock()
	defer g.Unlock()
	return g.State
}

func (g *GroupBase) SetCaptain(uid string) {
	for key, role := range g.Roles {
		if role == GroupRoleCaptain {
			g.Roles[key] = GroupRoleMember
		}
	}
	g.Roles[uid] = GroupRoleCaptain
}

func (g *GroupBase) CheckState(valids ...GroupState) error {
	for _, vs := range valids {
		if g.State == vs {
			return nil
		}
	}

	switch g.State {
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
		Captain:     g.GetCaptain(),
		GameMode:    g.GameMode,
		ModeVersion: g.ModeVersion,
		// TODO: other info
	}
}

func (g *GroupBase) CanStartMatch() error {
	return nil
}

func (g *GroupBase) SetAllowNearbyJoin(allow bool) {
	g.Settings.NearbyJoinAllowed = allow
}

func (g *GroupBase) SetAllowRecentJoin(allow bool) {
	g.Settings.RecentJoinAllowed = allow
}

func (g *GroupBase) AllowNearbyJoin() bool {
	return g.Settings.NearbyJoinAllowed
}

func (g *GroupBase) AllowRecentJoin() bool {
	return g.Settings.RecentJoinAllowed
}

func (g *GroupBase) AddInviteRecord(inviteeUID string, nowUnix int64) {
	g.InviteRecords[inviteeUID] = nowUnix + g.Configs.InviteExpireSec
}

func (g *GroupBase) DelInviteRecord(inviteeUID string) {
	delete(g.InviteRecords, inviteeUID)
}

func (g *GroupBase) GetInviteRecords() map[string]int64 {
	return g.InviteRecords
}
func (g *GroupBase) GetInviteExpireTimeStamp(uid string) int64 {
	return g.InviteRecords[uid]
}
func (g *GroupBase) IsInviteExpired(uid string, nowUnix int64) bool {
	return g.InviteRecords[uid] <= nowUnix
}

func (g *GroupBase) Lock() {
	g.L.Lock()
}

func (g *GroupBase) Unlock() {
	g.L.Unlock()
}

func (g *GroupBase) GetStartMatchTimeSec() int64 {
	return g.StartMatchTimeSec
}

func (g *GroupBase) SetStartMatchTimeSec(sec int64) {
	g.StartMatchTimeSec = sec
}
