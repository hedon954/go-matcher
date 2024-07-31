package common

import (
	"sync"

	"github.com/hedon954/go-matcher/config"
	"github.com/hedon954/go-matcher/merr"
	"github.com/hedon954/go-matcher/pto"
)

type GroupBase struct {
	// 在 GroupBase 的内部方法中不要进行同步处理，统一交给外部方法调用
	sync.RWMutex
	ID                int64
	Name              string
	players           []Player
	state             GroupState
	OwnerUID          string
	NearbyJoinAllowed bool
	RecentJoinAllowed bool

	// Config 队伍的相关配置
	Config config.GroupConfig

	/**
	 * 邀请卡片
	 *  key 有 2 种情况：
	 *   1. 玩家 Uid
	 *   2. 邀请渠道：single, clan, share, channel
	 */
	InvitedPlayers map[string]bool

	// 邀请附近的人, key: Uid, value: 附近的人 uids
	NearbyInviteMap map[string][]string

	// 小组中未准备的玩家，key: Uid
	UnReadyUsers map[string]bool

	// 队伍中有自己情侣，key: Uid
	CoupleMap map[string]bool

	// 每次匹配的独立数据
	MatchID string

	// 角色
	Roles map[string]int // TODO: 定义角色枚举
}

func NewGroupBase(id int64, p Player, c config.GroupConfig) *GroupBase {
	b := &GroupBase{
		ID:              id,
		state:           GroupStateInvite,
		OwnerUID:        p.UID(),
		players:         make([]Player, 0),
		InvitedPlayers:  make(map[string]bool),
		NearbyInviteMap: make(map[string][]string),
		UnReadyUsers:    make(map[string]bool),
		CoupleMap:       make(map[string]bool),
		Roles:           make(map[string]int),
		Config:          c,
	}
	b.AddPlayer(p)
	return b
}

func (b *GroupBase) GroupID() int64 {
	return b.ID
}

func (b *GroupBase) Base() *GroupBase {
	return b
}

func (b *GroupBase) GetMatchStrategy() int {
	if len(b.players) == 0 {
		return 0
	}
	return b.players[0].Base().MatchStrategy
}

func (b *GroupBase) GetGameMode() int {
	if len(b.players) == 0 {
		return 0
	}
	return b.players[0].Base().GameMode
}

func (b *GroupBase) GetModeVersion() int {
	if len(b.players) == 0 {
		return 0
	}
	return b.players[0].Base().ModeVersion
}

func (b *GroupBase) GetPlayers() []Player {
	return b.players
}

func (b *GroupBase) UIDs() []string {
	players := b.GetPlayers()
	uids := make([]string, len(players))
	for i := 0; i < len(uids); i++ {
		uids[i] = players[i].UID()
	}
	return uids
}

func (b *GroupBase) PlayerExists(p Player) bool {
	for _, gP := range b.players {
		if gP.UID() == p.UID() {
			return true
		}
	}
	return false
}

func (b *GroupBase) PlayerCount() int {
	return len(b.players)
}

func (b *GroupBase) AddPlayer(p Player) {
	if b.PlayerExists(p) {
		return
	} else {
		b.players = append(b.players, p)
		if len(b.players) == 1 {
			b.OwnerUID = p.UID()
		}
	}
}

func (b *GroupBase) DelPlayer(uid string) (empty bool) {
	deleted := false
	for i, p := range b.players {
		if p.UID() == uid {
			deleted = true
			b.players = append(b.players[:i], b.players[i+1:]...)
			break
		}
	}

	if deleted && len(b.players) > 0 {
		b.OwnerUID = b.players[0].UID()
	}

	return len(b.players) == 0
}

func (b *GroupBase) GetState() GroupState {
	return b.state
}

func (b *GroupBase) SetState(state GroupState) {
	b.state = state
}

func (b *GroupBase) CheckState(validStates ...GroupState) error {
	for _, state := range validStates {
		if b.state == state {
			return nil
		}
	}
	switch b.state {
	case GroupStateInvite:
		return merr.ErrGroupNotInQueue
	case GroupStateQueuing:
		return merr.ErrInQueuing
	case GroupStateMatched:
		return merr.ErrAlreadyMatched
	default:
		return merr.ErrUnknownGroupState
	}
}

func (b *GroupBase) IsOwner(uid string) bool {
	return b.OwnerUID == uid
}

func (b *GroupBase) GetOwner() Player {
	for _, p := range b.players {
		if p.UID() == b.OwnerUID {
			return p
		}
	}
	return nil
}

func (b *GroupBase) IsFull() bool {
	return b.PlayerCount() >= b.Config.PlayerLimit
}

func (b *GroupBase) CheckInvite() error {
	// 队伍状态
	if err := b.CheckState(GroupStateInvite); err != nil {
		return err
	}
	// 队伍人数
	if b.IsFull() {
		return merr.ErrGroupFull
	}
	return nil
}

func (b *GroupBase) GetGroupUsers() []pto.GroupUser {
	// TODO: get user groups
	return make([]pto.GroupUser, 0)
}

func (b *GroupBase) AddInvitedPlayer(inviteeUID string) {
	b.InvitedPlayers[inviteeUID] = true
}

func (b *GroupBase) CheckHandleInviteExpired(friendUID string, srcType pto.InvitationSrcType) error {
	// TODO: 优化
	switch srcType {
	case pto.InvitationSrcSingleChat:
		if _, ok := b.InvitedPlayers[string(SrcSingleChat)]; !ok {
			return merr.ErrInviteExpired
		}
	case pto.InvitationSrcClanRank:
		if _, ok := b.InvitedPlayers[string(SrcClanChat)]; !ok {
			return merr.ErrInviteExpired
		}
	case pto.InvitationSrcClanRace:
		if _, ok := b.InvitedPlayers[string(SrcClanChat)]; !ok {
			return merr.ErrInviteExpired
		}
	case pto.InvitationSrcShare:
		if _, ok := b.InvitedPlayers[string(SrcShare)]; !ok {
			return merr.ErrInviteExpired
		}
	case pto.InvitationSrcChannel:
		if _, ok := b.InvitedPlayers[string(SrcChannelChat)]; !ok {
			return merr.ErrInviteExpired
		}
	default:
		if _, ok := b.InvitedPlayers[friendUID]; !ok {
			return merr.ErrInviteExpired
		}
	}
	return nil
}

func (b *GroupBase) AddNearby(uid1, uid2 string) {
	b.NearbyInviteMap[uid1] = append(b.NearbyInviteMap[uid1], uid2)
	b.NearbyInviteMap[uid2] = append(b.NearbyInviteMap[uid2], uid1)
}

func (b *GroupBase) MarkSrcShare() {
	b.InvitedPlayers[string(SrcShare)] = true
}

func (b *GroupBase) SetCoupleInfo() {
	for _, p := range b.players {
		if _, ok := b.CoupleMap[p.UID()]; ok {
			p.SetWithCouple(true)
		}
	}
}

func (b *GroupBase) InitUnReadyMap() {
	// 需要准备的游戏模式重写该方法，默认不需要准备
	b.UnReadyUsers = make(map[string]bool)
}
