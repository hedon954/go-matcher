package glicko2

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	glicko "github.com/zelenin/go-glicko2"
)

func newQueue() *Queue {
	roomChan := make(chan Room, 128)
	q, _ := NewQueue("testQueue", roomChan, GetQueueArgs, NewTeam, NewRoom,
		NewRoomWithAi, zeroNowTime)
	return q
}

func zeroNowTime() int64 {
	return 0
}

func newPlayer(id string) *PlayerMock {
	return NewPlayer(id, false, 0, Args{})
}

func newPlayerWithMMR(id string, mmr float64) *PlayerMock {
	return NewPlayer(id, false, 0, Args{MMR: mmr})
}

func newGroupWithMMR(id int, playerCount int, mmr float64) Group {
	players := make([]*PlayerMock, playerCount)
	for i := 0; i < playerCount; i++ {
		players[i] = newPlayerWithMMR(cast.ToString(i*100+i), mmr)
	}
	g := NewGroup(cast.ToString(id), players)
	g.SetState(GroupStateQueuing)
	return g
}

func TestQueue_clearTmp(t *testing.T) {
	q := newQueue()
	groups := q.clearTmp()
	assert.Equal(t, 0, len(groups))
	g1 := NewGroup("1", []*PlayerMock{
		newPlayer("1"), newPlayer("2"),
	})
	g1.SetState(GroupStateQueuing)
	g2 := NewGroup("2", []*PlayerMock{
		newPlayer("1"), newPlayer("2"),
	})
	g2.SetState(GroupStateQueuing)
	t1 := NewTeam(g1)
	t2 := NewTeam(g2)
	q.TmpTeam = append(q.TmpTeam, t1)
	q.TmpTeam = append(q.TmpTeam, t2)
	groups = q.clearTmp()
	assert.Equal(t, 2, len(groups))
	assert.Equal(t, 0, len(q.TmpTeam))
	assert.Equal(t, 0, len(q.TmpRoom))
	assert.Equal(t, "1", groups[0].GetID())
	assert.Equal(t, "2", groups[1].GetID())
	assert.Equal(t, GroupStateQueuing, groups[0].GetState())
	assert.Equal(t, GroupStateQueuing, groups[1].GetState())
	t1 = NewTeam(g1)
	t2 = NewTeam(g2)
	r1 := NewRoomWithAi(t1)
	r2 := NewRoomWithAi(t2)
	assert.Equal(t, true, r1.HasAi())
	assert.Equal(t, true, r2.HasAi())
	q.TmpRoom = append(q.TmpRoom, r1)
	q.TmpRoom = append(q.TmpRoom, r2)
	groups = q.clearTmp()
	assert.Equal(t, 6, len(groups))
	assert.Equal(t, 0, len(q.TmpTeam))
	assert.Equal(t, 0, len(q.TmpRoom))
	assert.Equal(t, "1", groups[0].GetID())
	assert.Equal(t, "2", groups[3].GetID())
	assert.Equal(t, GroupStateQueuing, groups[0].GetState())
	assert.Equal(t, GroupStateQueuing, groups[3].GetState())
}

func TestQueue_AddGroups_AllGroups_GetAndClearGroups(t *testing.T) {
	q := newQueue()
	groups := q.GetAndClearGroups()
	assert.Equal(t, 0, len(groups))
	g1 := NewGroup("1", []*PlayerMock{
		newPlayer("1"), newPlayer("2"),
	})
	g1.SetState(GroupStateQueuing)
	g2 := NewGroup("2", []*PlayerMock{
		newPlayer("1"), newPlayer("2"),
	})
	q.AddGroups(g1, g2)
	groups = q.AllGroups()
	assert.Equal(t, 2, len(groups))
	groups = q.GetAndClearGroups()
	assert.Equal(t, 1, len(groups))
	assert.Equal(t, "1", groups[0].GetID())
	groups = q.GetAndClearGroups()
	assert.Equal(t, 0, len(groups))
}

func TestQueue_getMatchRange(t *testing.T) {
	q := newQueue()
	q.QueueArgs = &QueueArgs{}
	// 没有配置，则拿默认的
	mr := q.getMatchRange(q.nowUnixFunc()-1, q.nowUnixFunc()-1)
	assert.Equal(t, defaultMatchRange, mr)
	q.QueueArgs = q.getQueueArgs()
	// 有配置的，根据配置的拿
	mr = q.getMatchRange(q.nowUnixFunc()-0, q.nowUnixFunc()-0)
	assert.Equal(t, MatchRange{
		MaxMatchSec:   15,
		MMRGapPercent: 10,
		CanJoinTeam:   false,
		StarGap:       0,
	}, mr)
	mr = q.getMatchRange(q.nowUnixFunc()-0, q.nowUnixFunc()-16)
	assert.Equal(t, MatchRange{
		MaxMatchSec:   15,
		MMRGapPercent: 10,
		CanJoinTeam:   false,
		StarGap:       0,
	}, mr)
	mr = q.getMatchRange(q.nowUnixFunc()-20, q.nowUnixFunc()-16)
	assert.Equal(t, MatchRange{
		MaxMatchSec:   30,
		MMRGapPercent: 20,
		CanJoinTeam:   false,
		StarGap:       0,
	}, mr)
	mr = q.getMatchRange(q.nowUnixFunc()-0, q.nowUnixFunc()-1000)
	assert.Equal(t, MatchRange{
		MaxMatchSec:   15,
		MMRGapPercent: 10,
		CanJoinTeam:   false,
		StarGap:       0,
	}, mr)
	mr = q.getMatchRange(q.nowUnixFunc()-1000, q.nowUnixFunc()-1000)
	assert.Equal(t, MatchRange{
		MaxMatchSec:   90,
		MMRGapPercent: 0,
		CanJoinTeam:   true,
		StarGap:       0,
	}, mr)
}

func TestQueue_buildNewTeams(t *testing.T) {
	q := newQueue()
	g1 := NewGroup("1", []*PlayerMock{newPlayer("1"), newPlayer("2"), newPlayer("3")})
	g1.SetState(GroupStateQueuing)
	g2 := NewGroup("2", []*PlayerMock{newPlayer("4"), newPlayer("5")})
	g2.SetState(GroupStateQueuing)
	g3 := NewGroup("3", []*PlayerMock{newPlayer("6"), newPlayer("7")})
	g3.SetState(GroupStateQueuing)
	g4 := NewGroup("4", []*PlayerMock{newPlayer("8"), newPlayer("9"), newPlayer("10")})
	g4.SetState(GroupStateQueuing)
	g5 := NewGroup("5", []*PlayerMock{newPlayer("11")})
	g5.SetState(GroupStateQueuing)
	q.AddGroups(g1, g2, g3, g4, g5)
	groups := q.GetAndClearGroups()
	sort.Slice(groups, func(i, j int) bool {
		return cast.ToInt(groups[i].GetID()) < cast.ToInt(groups[j].GetID())
	})
	assert.Equal(t, 5, len(groups))
	groups = q.buildNewTeams(groups)
	groups = q.GetAndClearGroups()
	q.AddGroups(groups...)
	sort.Slice(groups, func(i, j int) bool {
		return cast.ToInt(groups[i].GetID()) < cast.ToInt(groups[j].GetID())
	})
	assert.Equal(t, 0, len(groups))
	assert.Equal(t, 2, len(q.FullTeam))
	assert.Equal(t, 1, len(q.TmpTeam))
}

func TestQueue_buildNewRooms(t *testing.T) {
	q := newQueue()
	g1 := NewGroup("1", []*PlayerMock{newPlayer("1"), newPlayer("2"), newPlayer("3")})
	g1.SetState(GroupStateQueuing)
	g2 := NewGroup("2", []*PlayerMock{newPlayer("4"), newPlayer("5")})
	g2.SetState(GroupStateQueuing)
	g3 := NewGroup("3", []*PlayerMock{newPlayer("6"), newPlayer("7")})
	g3.SetState(GroupStateQueuing)
	g4 := NewGroup("4", []*PlayerMock{newPlayer("8"), newPlayer("9"), newPlayer("10")})
	g4.SetState(GroupStateQueuing)
	g5 := NewGroup("5", []*PlayerMock{newPlayer("11")})
	g5.SetState(GroupStateQueuing)
	q.AddGroups(g1, g2, g3, g4, g5)
	groups := q.GetAndClearGroups()
	sort.Slice(groups, func(i, j int) bool {
		return cast.ToInt(groups[i].GetID()) < cast.ToInt(groups[j].GetID())
	})
	assert.Equal(t, 5, len(groups))
	groups = q.buildNewTeams(groups)

	// 构建房间，只有人数 >=5 的 team 才可以加入
	q.buildNewRooms()
	assert.Equal(t, 1, len(q.TmpRoom))
	assert.Equal(t, 2, len(q.TmpRoom[0].GetTeams()))
	assert.Equal(t, 1, len(q.TmpTeam))
	assert.Equal(t, 1, len(q.TmpTeam[0].GetGroups()))
	assert.Equal(t, "5", q.TmpTeam[0].GetGroups()[0].GetID())
}

func TestQueue_buildNewTeam_fillTmpTeam_buildNewRoom(t *testing.T) {
	q := newQueue()
	g1 := NewGroup("1", []*PlayerMock{newPlayer("1"), newPlayer("2"), newPlayer("3")})
	g1.SetState(GroupStateQueuing)
	g2 := NewGroup("2", []*PlayerMock{newPlayer("4"), newPlayer("5")})
	g2.SetState(GroupStateQueuing)
	g3 := NewGroup("3", []*PlayerMock{newPlayer("6"), newPlayer("7")})
	g3.SetState(GroupStateQueuing)
	g4 := NewGroup("4", []*PlayerMock{newPlayer("8"), newPlayer("9"), newPlayer("10")})
	g4.SetState(GroupStateQueuing)
	g5 := NewGroup("5", []*PlayerMock{newPlayer("11")})
	g5.SetState(GroupStateQueuing)
	q.AddGroups(g1, g2, g3, g4, g5)
	groups := q.GetAndClearGroups()
	sort.Slice(groups, func(i, j int) bool {
		return cast.ToInt(groups[i].GetID()) < cast.ToInt(groups[j].GetID())
	})

	assert.Equal(t, 5, len(groups))
	groups = q.buildNewTeams(groups)

	// 构建房间，只有人数 >=5 的 team 才可以加入
	q.buildNewRooms()
	assert.Equal(t, 1, len(q.TmpRoom))
	assert.Equal(t, 2, len(q.TmpRoom[0].GetTeams()))
	assert.Equal(t, 1, len(q.TmpTeam))
	assert.Equal(t, 1, len(q.TmpTeam[0].GetGroups()))
	assert.Equal(t, "5", q.TmpTeam[0].GetGroups()[0].GetID())
}

// 1. 第1个 group 直接进 team
// 2. 尽可能找 mmr 相近的 group 组成一个 team
// 3. 只有符合 mmr 匹配范围的 group，才可以组成一个 team
// 4. 随着匹配时间的增长，不断打开 group mmr 的匹配范围
func TestQueue_findGroupForTeam(t *testing.T) {
	var found bool
	q := newQueue()

	// 第 1 个 group 直接进
	t1 := &TeamMock{groups: make(map[string]Group)}
	g1 := newGroupWithMMR(1, 3, 10)
	g1.SetState(GroupStateQueuing)
	_, found = q.findGroupForTeamBinary(t1, []Group{g1})
	assert.Equal(t, true, found)
	assert.Equal(t, 1, len(t1.GetGroups()))

	/**
		{
			MaxMatchSec:   15,
			MMRGapPercent: 10,
			CanJoinTeam:   false,
		},
		{
			MaxMatchSec:   30,
			MMRGapPercent: 20,
			CanJoinTeam:   false,
		},
		{
			MaxMatchSec:   60,
			MMRGapPercent: 30,
			CanJoinTeam:   true,
		},
		{
			MaxMatchSec:   90,
			MMRGapPercent: 0,
			CanJoinTeam:   true,
		},

		// 匹配时间超过 60s 接受 AI
	},
	*/

	// g1 和 g2 mmr 相差 70，匹配时间在 15s 内，不匹配，无法组成一队
	g2 := newGroupWithMMR(2, 2, 100)
	groups := []Group{g2}
	groups, found = q.findGroupForTeamBinary(t1, groups)
	assert.Equal(t, false, found)
	assert.Equal(t, 1, len(t1.GetGroups()))
	assert.Equal(t, GroupStateQueuing, g2.GetState())
	assert.Equal(t, 1, len(groups))

	// g1 和 g3 mmr 相差 3，匹配时间在 15 秒内，< 30%，可以组成一队
	g3 := newGroupWithMMR(3, 2, 11)
	groups = append(groups, g3)
	groups, found = q.findGroupForTeamBinary(t1, groups)
	assert.Equal(t, 2, len(t1.GetGroups()))

	// g2(100) 和 g4(130) mmr 相差 30，匹配时间 25s，只能接受 20% 的差距，无法组成一队
	t2 := &TeamMock{groups: make(map[string]Group)}
	groups, found = q.findGroupForTeamBinary(t2, groups)
	assert.Equal(t, 1, len(t2.GetGroups()))
	g4 := newGroupWithMMR(4, 3, 130)
	groups = append(groups, g4)
	g2.SetStartMatchTimeSec(q.nowUnixFunc() - 25)
	g4.SetStartMatchTimeSec(q.nowUnixFunc() - 55)
	groups, found = q.findGroupForTeamBinary(t2, groups)
	assert.Equal(t, false, found)
	assert.Equal(t, 1, len(t2.groups))
	assert.Equal(t, GroupStateQueuing, g4.GetState())

	// g2 和 g4 mmr 相差 30，匹配时间 45s，可以接受 30% 的差距，
	// 且两队时间都在 60s 内，接受 AI，
	// 可以组成一队
	g2.SetStartMatchTimeSec(q.nowUnixFunc() - 45)
	groups, found = q.findGroupForTeamBinary(t2, groups)
	assert.Equal(t, true, found)
	assert.Equal(t, 2, len(t2.groups))

	// g5 g6 相差很大，但是匹配很久了，可以直接组队
	g5 := newGroupWithMMR(5, 4, 10000)
	g5.SetStartMatchTimeSec(q.nowUnixFunc() - 10000)
	g6 := newGroupWithMMR(6, 1, 1)
	g6.SetStartMatchTimeSec(q.nowUnixFunc() - 50000)
	groups = append(groups, g5, g6)
	t3 := NewTeam(g5)
	groups, found = q.findGroupForTeamBinary(t3, groups)
	assert.Equal(t, true, found)
	assert.Equal(t, 2, len(t3.GetGroups()))

	// g7,g8,g9,g10，g7 和 g9 更近，应该选 g9
	g7 := newGroupWithMMR(7, 3, 100)
	g8 := newGroupWithMMR(8, 2, 110)
	g9 := newGroupWithMMR(9, 2, 105)
	g10 := newGroupWithMMR(10, 2, 109)
	groups = append(groups, g7, g8, g9, g10)
	t4 := NewTeam(g7)
	groups, found = q.findGroupForTeamBinary(t4, groups)
	assert.Equal(t, true, found)
	assert.Equal(t, 2, len(t4.GetGroups()))
	assert.Equal(t, g9, t4.GetGroups()[1])
	g8.SetState(GroupStateUnready)
	g10.SetState(GroupStateUnready)

	// g11,g12,g13
	g11 := newGroupWithMMR(11, 3, 100)
	g12 := newGroupWithMMR(12, 1, 93)
	g13 := newGroupWithMMR(13, 1, 108)
	groups = append(groups, g11, g12, g13)
	t5 := NewTeam(g11)
	// g12 符合要求，进队
	groups, found = q.findGroupForTeamBinary(t5, groups)
	assert.Equal(t, true, found)
	assert.Equal(t, 2, len(t5.GetGroups()))
	assert.Equal(t, GroupStateQueuing, g13.GetState())
	// g13 符合 g11 的要求，但是不符合 g12 的要求，无法进队
	groups, found = q.findGroupForTeamBinary(t5, groups)
	assert.Equal(t, false, found)
	assert.Equal(t, 2, len(t5.GetGroups()))
	assert.Equal(t, GroupStateQueuing, g13.GetState())
}

// 测试 findGroupForTeam 在遍历情况下的快速退出
func TestQueue_findGroupForTeam_rangeSearch_quickReturn(t *testing.T) {
	found := false
	q := newQueue()
	g1 := newGroupWithMMR(1, 2, 1000)
	t1 := NewTeam(g1)
	q.TmpTeam = append(q.TmpTeam, t1)
	g2 := newGroupWithMMR(2, 1, 910)
	g3 := newGroupWithMMR(3, 1, 930)
	g4 := newGroupWithMMR(4, 1, 970)
	g5 := newGroupWithMMR(5, 1, 1020)
	g6 := newGroupWithMMR(6, 1, 1040)
	g7 := newGroupWithMMR(7, 1, 1050)
	g8 := newGroupWithMMR(8, 1, 1090)
	q.AddGroups(g2, g3, g4, g5, g6, g7, g8)
	groups := q.SortedGroups()
	groups, found = q.findGroupForTeamRange(t1, groups)
	assert.Equal(t, true, found)
	assert.Equal(t, 2, len(t1.GetGroups()))
	assert.Equal(t, g5, t1.GetGroups()[1])
}

// 1. 只有人数为 5 的 team 可以进入 romm
// 2. 第 1 个 team 直接进入 room
// 3. 尽可能找相近的 team 组成一个 room
// 4. 只有符合条件的 mmr 可以组成一个 room
// 5. 随时匹配时间增长，逐渐放开 mmr 的匹配范围和是否允许和车队组队
func TestQueue_findTeamForRoom(t *testing.T) {
	q := newQueue()

	// 1. 只有人数为 5 的 team 可以进入 romm
	g1 := newGroupWithMMR(1, 4, 100)
	q.TmpTeam = append(q.TmpTeam, NewTeam(g1))
	r1 := &RoomMock{teams: make([]Team, 0)}
	q.findTeamForRoomBinary(r1)
	assert.Equal(t, 0, len(r1.teams))
	assert.Equal(t, 1, len(q.TmpTeam))

	// 2. 第 1 个 team 直接进入 room
	g2 := newGroupWithMMR(2, 5, 100)
	t2 := NewTeam(g2)
	q.FullTeam = append(q.FullTeam, t2)
	q.findTeamForRoomBinary(r1)
	assert.Equal(t, 1, len(r1.teams))
	assert.Equal(t, r1.teams[0], t2)
	assert.Equal(t, 1, len(q.TmpTeam))

	// 3. 尽可能找相近的 team 组成一个 room
	g3 := newGroupWithMMR(3, 5, 101)
	t3 := NewTeam(g3)
	g4 := newGroupWithMMR(4, 5, 104)
	t4 := NewTeam(g4)
	g5 := newGroupWithMMR(5, 5, 103)
	t5 := NewTeam(g5)
	g6 := newGroupWithMMR(6, 5, 102)
	t6 := NewTeam(g6)
	q.FullTeam = append(q.FullTeam, t3, t4, t5, t6)
	q.findTeamForRoomBinary(r1)
	// 先 t3
	assert.Equal(t, 2, len(r1.teams))
	assert.Equal(t, r1.teams[1], t3)
	// 再 t6
	q.findTeamForRoomBinary(r1)
	assert.Equal(t, 3, len(r1.teams))
	assert.Equal(t, r1.teams[2], t6)

	// 4. 只有符合条件的 mmr 可以组成一个 room
	r2 := &RoomMock{teams: make([]Team, 0)}
	q.findTeamForRoomBinary(r2) // 目前 FullTeam [t4,t5]
	// t4 直接进
	assert.Equal(t, 1, len(r2.teams))
	assert.Equal(t, t4, r2.teams[0])
	// t5 符合条件，也进
	q.findTeamForRoomBinary(r2)
	assert.Equal(t, 2, len(r2.teams))
	assert.Equal(t, t5, r2.teams[1])
	g7 := newGroupWithMMR(7, 5, 120)
	t7 := NewTeam(g7)
	q.FullTeam = append(q.FullTeam, t7)
	// t7 不符合条件，进不来
	q.findTeamForRoomBinary(r2)
	assert.Equal(t, 2, len(r2.teams))
	assert.Equal(t, 1, len(q.FullTeam))
	// 5. 随时匹配时间增长，逐渐放开 mmr 的匹配范围和是否允许和车队组队
	t4.(*TeamMock).StartMatchTimeSec = q.nowUnixFunc() - 1000
	t5.(*TeamMock).StartMatchTimeSec = q.nowUnixFunc() - 1000
	t7.(*TeamMock).StartMatchTimeSec = q.nowUnixFunc() - 1000
	q.findTeamForRoomBinary(r2)
	assert.Equal(t, 3, len(r2.teams))
	assert.Equal(t, t7, r2.teams[2])
	// 6. t8,t9 不是车队,t10是车队
	g8 := newGroupWithMMR(8, 3, 100)
	g9 := newGroupWithMMR(9, 2, 100)
	t8 := NewTeam(g8)
	q.findGroupForTeamBinary(t8, []Group{g9})
	assert.Equal(t, 2, len(t8.GetGroups()))
	g10 := newGroupWithMMR(10, 3, 105)
	g11 := newGroupWithMMR(11, 2, 105)
	t9 := NewTeam(g10)
	q.findGroupForTeamBinary(t9, []Group{g11})
	assert.Equal(t, 2, len(t9.GetGroups()))
	g12 := newGroupWithMMR(12, 5, 104)
	t10 := NewTeam(g12)
	q.FullTeam = append(q.FullTeam, t8, t9, t10)
	r3 := &RoomMock{teams: make([]Team, 0)}
	// t8 直接进
	q.findTeamForRoomBinary(r3)
	assert.Equal(t, 1, len(r3.teams))
	assert.Equal(t, t8, r3.teams[0])
	// t10 比 t9 更接近 t8，但是因为是车队，所以不符合
	q.findTeamForRoomBinary(r3)
	assert.Equal(t, 2, len(r3.teams))
	assert.Equal(t, t9, r3.teams[1])
	// t10 进不去
	q.findTeamForRoomBinary(r3)
	assert.Equal(t, 2, len(r3.teams))
	// 随着匹配时间加成，t8,t9 支持匹配车队
	t8.(*TeamMock).StartMatchTimeSec = q.nowUnixFunc() - 1000
	t9.(*TeamMock).StartMatchTimeSec = q.nowUnixFunc() - 1000
	t10.(*TeamMock).StartMatchTimeSec = q.nowUnixFunc() - 1000
	q.findTeamForRoomBinary(r3)
	assert.Equal(t, 3, len(r3.teams))
}

// ==================================================
//                  下面的接口的 mock
// ==================================================

const (
	// 车队在专属队列中的匹配时长
	NormalTeamWaitTimeSec     int64 = 5
	UnfriendlyTeamWaitTimeSec int64 = 10
	MaliciousTeamWaitTimeSec  int64 = 15

	TeamPlayerLimit = 5 // 阵营总人数
	RoomTeamLimit   = 3 // 房间总阵营数
)

func GetQueueArgs() *QueueArgs {
	return &QueueArgs{
		TeamPlayerLimit:           TeamPlayerLimit,
		RoomTeamLimit:             RoomTeamLimit,
		NormalTeamWaitTimeSec:     NormalTeamWaitTimeSec,
		UnfriendlyTeamWaitTimeSec: UnfriendlyTeamWaitTimeSec,
		MaliciousTeamWaitTimeSec:  MaliciousTeamWaitTimeSec,
		MatchRanges: []MatchRange{
			{
				MaxMatchSec:   15,
				MMRGapPercent: 10,
				CanJoinTeam:   false,
				StarGap:       0,
			},
			{
				MaxMatchSec:   30,
				MMRGapPercent: 20,
				CanJoinTeam:   false,
				StarGap:       0,
			},
			{
				MaxMatchSec:   60,
				MMRGapPercent: 30,
				CanJoinTeam:   true,
				StarGap:       0,
			},
			{
				MaxMatchSec:   90,
				MMRGapPercent: 0,
				CanJoinTeam:   true,
				StarGap:       0,
			},
		},
	}
}

type PlayerMock struct {
	sync.RWMutex `json:"-"`

	ID      string `json:"ID"`
	isAi    bool
	aiLevel int64

	MMR float64 `json:"-"`
	RD  float64 `json:"-"`
	V   float64 `json:"-"`

	rank int
	star int

	startMatchTime  int64
	finishMatchTime int64

	*glicko.Player `json:"-"`
}

func NewPlayer(id string, isAi bool, aiLevel int64, args Args) *PlayerMock {
	/**
	TODO:
	算法刚启动的时候，会手动配置不同段位的补偿分数，把各个段位的分数人工区分初始化玩家的 MMR，rd 和 v，

		赛季重置也会初始化玩家的相关分数，重置方式如下；
	1. 初始评分（MMR），转换为上赛季评分*70%，最低分1000，最高分7500
	2. RD，根据当前赛季和历史最高赛季的星星数差距决定，最高700，最低为0
	3. 波动率，初始为0.06
	*/
	return &PlayerMock{
		RWMutex: sync.RWMutex{},
		ID:      id,
		isAi:    isAi,
		aiLevel: aiLevel,
		MMR:     args.MMR,
		RD:      args.RD,
		V:       args.V,
		Player:  glicko.NewPlayer(glicko.NewRating(args.MMR, args.RD, args.V)),
	}
}

func (p *PlayerMock) GlickoPlayer() *glicko.Player {
	return p.Player
}

func (p *PlayerMock) IsAi() bool {
	return p.isAi
}

func (p *PlayerMock) GetAiLevel() int64 {
	return p.aiLevel
}

func (p *PlayerMock) GetStartMatchTimeSec() int64 {
	return p.startMatchTime
}

func (p *PlayerMock) SetStartMatchTimeSec(t int64) {
	p.startMatchTime = t
}

func (p *PlayerMock) GetFinishMatchTimeSec() int64 {
	return p.finishMatchTime
}

func (p *PlayerMock) SetFinishMatchTimeSec(t int64) {
	p.finishMatchTime = t
}

func (p *PlayerMock) GetID() string {
	return p.ID
}

func (p *PlayerMock) GetMMR() float64 {
	p.RLock()
	defer p.RUnlock()
	return p.MMR
}

func (p *PlayerMock) GetStar() int {
	return p.star
}

func (p *PlayerMock) GetArgs() *Args {
	p.RLock()
	defer p.RUnlock()
	/**
	TODO:
	赛季初始 5 局过后，MMR 分数开始生效；生效前沿用上赛季分数进行匹配。
	新玩家和回流玩家会进入保护期，分数不计算，对局全部按照最低分进行匹配，直到完成10次团战对局，开始使用真实分数进行匹配。
	*/
	return &Args{
		MMR: p.MMR,
		RD:  p.RD,
		V:   p.V,
	}
}

func (p *PlayerMock) SetArgs(args *Args) error {
	if args == nil {
		return errors.New("args is nil")
	}
	p.Lock()
	defer p.Unlock()
	p.V = args.V
	p.MMR = args.MMR
	p.RD = args.RD
	p.Player = glicko.NewPlayer(glicko.NewRating(args.MMR, args.RD, args.V))
	return nil
}

func (p *PlayerMock) SetStar(star int) {
	p.star = star
}

func (p *PlayerMock) GetRank() int {
	return p.rank
}

func (p *PlayerMock) SetRank(rank int) {
	p.rank = rank
}

const (
	// 车队方差阈值
	MaliciousTeamVarianceMin  = 100000
	UnfriendlyTeamVarianceMin = 1000
)

type GroupMock struct {
	sync.RWMutex

	ID         string              `json:"id"`
	State      GroupState          `json:"state"`
	PlayersMap map[string]struct{} `json:"players_map"`
	Players    []*PlayerMock       `json:"players"`

	modeVersion int
	platform    int

	startMatchTimeSec int64
}

func (g *GroupMock) IsNewer() bool {
	return false
}

func NewGroup(id string, players []*PlayerMock) *GroupMock {
	g := &GroupMock{
		RWMutex:    sync.RWMutex{},
		ID:         id,
		State:      GroupStateUnready,
		PlayersMap: make(map[string]struct{}),
		Players:    players,
	}
	for _, p := range g.Players {
		g.PlayersMap[p.GetID()] = struct{}{}
		g.startMatchTimeSec = p.GetStartMatchTimeSec()
	}
	return g
}

func (g *GroupMock) GetID() string {
	return g.ID
}

func (g *GroupMock) QueueKey() string {
	return fmt.Sprintf("%d-%d", g.modeVersion, g.platform)
}

func (g *GroupMock) GetState() GroupState {
	g.RLock()
	defer g.RUnlock()

	return g.State
}

func (g *GroupMock) SetState(state GroupState) {
	g.Lock()
	defer g.Unlock()

	g.State = state
}

func (g *GroupMock) GetModeVersion() int {
	return g.modeVersion
}

func (g *GroupMock) GetPlatform() int {
	return g.platform
}

func (g *GroupMock) PlayerCount() int {
	g.RLock()
	defer g.RUnlock()
	return len(g.Players)
}

func (g *GroupMock) GetPlayers() []Player {
	g.RLock()
	defer g.RUnlock()
	res := make([]Player, len(g.Players))
	for i := 0; i < len(res); i++ {
		res[i] = g.Players[i]
	}
	return res
}

func (g *GroupMock) ForceCancelMatch(reason string, waitSec int64) {

}

func (g *GroupMock) AddPlayers(players ...Player) {
	g.Lock()
	defer g.Unlock()

	for _, p := range players {
		_, ok := g.PlayersMap[p.GetID()]
		if ok {
			continue
		}
		g.PlayersMap[p.GetID()] = struct{}{}
		g.Players = append(g.Players, p.(*PlayerMock))
	}
}

func (g *GroupMock) RemovePlayers(players ...Player) {
	g.Lock()
	defer g.Unlock()

	for _, p := range players {
		for i, gp := range g.Players {
			if gp == p {
				g.Players = append(g.Players[:i], g.Players[i+1:]...)
				delete(g.PlayersMap, p.GetID())
			}
		}
	}
}

// AverageMMR 算出队伍的平均 MMR
func (g *GroupMock) AverageMMR() float64 {
	total := 0.0
	for _, player := range g.Players {
		total += player.GetMMR()
	}
	return total / float64(len(g.Players))
}

// MMR 算出队伍的最大的 MMR
func (g *GroupMock) BiggestMMR() float64 {
	mmr := 0.0
	for _, p := range g.Players {
		pMMR := p.GetMMR()
		if pMMR > mmr {
			mmr = pMMR
		}
	}
	return mmr
}

// MMR 算出队伍的 MMR
func (g *GroupMock) GetMMR() float64 {
	teamType := g.Type()
	switch teamType {
	case GroupTypeUnfriendlyTeam:
		mmr := g.AverageMMR() * 1.5
		bMmr := g.BiggestMMR()
		if mmr > bMmr {
			mmr = bMmr
		}
		return mmr
	case GroupTypeMaliciousTeam:
		return g.BiggestMMR()
	default:
		return g.AverageMMR()
	}
}

// Rank 队伍段位要弄平均值替代
func (g *GroupMock) GetStar() int {
	if len(g.Players) == 0 {
		return 0
	}
	rank := 0
	for _, p := range g.Players {
		rank += p.GetStar()
	}
	return rank / len(g.Players)
}

// Group 算出队伍的 MMR 方差
func (g *GroupMock) MMRVariance() float64 {
	data := stats.Float64Data{}
	for _, p := range g.Players {
		data = append(data, p.GetMMR())
	}
	variance, _ := stats.Variance(data)
	return variance
}

// Type 确定车队类型
func (g *GroupMock) Type() GroupType {
	if len(g.Players) != 5 {
		return GroupTypeNotTeam
	}
	variance := g.MMRVariance()
	if variance >= MaliciousTeamVarianceMin {
		return GroupTypeMaliciousTeam
	} else if variance >= UnfriendlyTeamVarianceMin {
		return GroupTypeUnfriendlyTeam
	} else {
		return GroupTypeNormalTeam
	}
}

func (g *GroupMock) CanFillAi() bool {
	now := zeroNowTime()
	if now-g.GetStartMatchTimeSec() > 60 {
		return true
	}
	return false
}

func (g *GroupMock) Print() {
	fmt.Printf("\t\t%s\t\t\t%d\t\t%.2f\t\t%.2f\t\t%ds\t\t\n", g.GetID(), len(g.Players), g.GetMMR(), g.AverageMMR(),
		time.Now().Unix()-g.GetStartMatchTimeSec())
}

func (g *GroupMock) GetFinishMatchTimeSec() int64 {
	if len(g.Players) == 0 {
		return 0
	}
	return g.Players[0].GetFinishMatchTimeSec()
}

func (g *GroupMock) SetFinishMatchTimeSec(t int64) {
	for _, p := range g.Players {
		p.SetFinishMatchTimeSec(t)
	}
}

func (g *GroupMock) GetStartMatchTimeSec() int64 {
	return g.startMatchTimeSec
}

func (g *GroupMock) SetStartMatchTimeSec(t int64) {
	g.startMatchTimeSec = t
	for _, p := range g.Players {
		p.SetStartMatchTimeSec(t)
	}
}

type RoomMock struct {
	id              int64
	teams           []Team
	StartMatchTime  int64
	FinishMatchTime int64
}

func NewRoom(team Team) Room {
	r := &RoomMock{
		teams: make([]Team, 0, 3),
	}
	r.AddTeam(team)
	return r
}

func NewRoomWithAi(team Team) Room {
	newRoom := NewRoom(team)
	ai1G := NewGroup("isAi-group-0", nil)
	for i := 0; i < TeamPlayerLimit; i++ {
		ai1G.AddPlayers(NewPlayer("isAi-player-0-"+strconv.Itoa(i), true, int64(i+1), Args{}))
	}
	ai1G.SetState(GroupStateQueuing)
	aiT1 := NewTeam(ai1G)
	ai2G := NewGroup("isAi-group-1", nil)
	for i := 0; i < TeamPlayerLimit; i++ {
		ai2G.AddPlayers(NewPlayer("isAi-player-1-"+strconv.Itoa(i), true, int64(i+1), Args{}))
	}
	ai2G.SetState(GroupStateQueuing)
	aiT2 := NewTeam(ai2G)
	newRoom.AddTeam(aiT1)
	newRoom.AddTeam(aiT2)
	return newRoom
}

func (r *RoomMock) GetID() int64 {
	return r.id
}

func (r *RoomMock) SetID(rid int64) {
	r.id = rid
}

func (r *RoomMock) GetTeams() []Team {
	return r.teams
}

func (r *RoomMock) AddTeam(t Team) {
	if t.PlayerCount() != 5 {
		fmt.Print()
	}
	r.teams = append(r.teams, t)
	tmst := t.GetStartMatchTimeSec()
	if tmst == 0 {
		return
	}
	if r.StartMatchTime == 0 || r.StartMatchTime > tmst {
		r.StartMatchTime = tmst
	}
}

func (r *RoomMock) RemoveTeam(t Team) {
	for i, rt := range r.teams {
		if rt == t {
			r.teams = append(r.teams[:i], r.teams[i+1:]...)
			break
		}
	}
	return
}

func (r *RoomMock) GetMMR() float64 {
	if len(r.teams) == 0 {
		return 0.0
	}
	mmr := 0.0
	for _, t := range r.teams {
		mmr += t.GetMMR()
	}
	return mmr / float64(len(r.teams))
}

func (r *RoomMock) GetStartMatchTimeSec() int64 {
	return r.StartMatchTime
}

func (r *RoomMock) GetFinishMatchTimeSec() int64 {
	return r.FinishMatchTime
}

func (r *RoomMock) PlayerCount() int {
	count := 0
	for _, t := range r.teams {
		count += t.PlayerCount()
	}
	return count
}

func (r *RoomMock) SetFinishMatchTimeSec(t int64) {
	for _, team := range r.teams {
		team.SetFinishMatchTimeSec(t)
	}
	r.FinishMatchTime = t
}

func (r *RoomMock) HasAi() bool {
	for _, t := range r.teams {
		if t.IsAi() {
			return true
		}
	}
	return false
}

type TeamMock struct {
	sync.RWMutex

	groups            map[string]Group
	StartMatchTimeSec int64
	rank              int
}

func (t *TeamMock) IsFull(teamPlayerLimit int) bool {
	return t.PlayerCount() >= teamPlayerLimit
}

func (t *TeamMock) IsNewer() bool {
	return false
}

func NewTeam(group Group) Team {
	t := &TeamMock{
		RWMutex: sync.RWMutex{},
		groups:  make(map[string]Group),
	}
	t.AddGroup(group)
	return t
}

func (t *TeamMock) GetModeVersion() int {
	return 0
}

func (t *TeamMock) GetGroups() []Group {
	t.RLock()
	defer t.RUnlock()
	res := make([]Group, len(t.groups))
	i := 0
	for _, g := range t.groups {
		res[i] = g
		i++
	}
	sort.Slice(res, func(i, j int) bool {
		return cast.ToInt(res[i].GetID()) < cast.ToInt(res[j].GetID())
	})
	return res
}

func (t *TeamMock) AddGroup(g Group) {
	t.Lock()
	defer t.Unlock()
	t.groups[g.GetID()] = g
	gmst := g.GetStartMatchTimeSec()
	if gmst == 0 {
		return
	}
	if t.StartMatchTimeSec == 0 || t.StartMatchTimeSec > gmst {
		t.StartMatchTimeSec = gmst
	}
}

func (t *TeamMock) RemoveGroup(groupId string) {
	t.Lock()
	defer t.Unlock()
	delete(t.groups, groupId)
}

func (t *TeamMock) PlayerCount() int {
	t.RLock()
	defer t.RUnlock()
	count := 0
	for _, group := range t.groups {
		count += group.PlayerCount()
	}
	return count
}

func (t *TeamMock) GetMMR() float64 {
	t.RLock()
	defer t.RUnlock()
	if len(t.groups) == 0 {
		return 0
	}
	total := 0.0
	for _, group := range t.groups {
		total += group.GetMMR()
	}
	return total / float64(len(t.groups))
}

func (t *TeamMock) GetStar() int {
	t.RLock()
	defer t.RUnlock()
	if len(t.groups) == 0 {
		return 0
	}
	rank := 0
	for _, g := range t.groups {
		rank += g.GetStar()
	}
	return rank / len(t.groups)
}

func (t *TeamMock) SetFinishMatchTimeSec(t2 int64) {
	t.RLock()
	defer t.RUnlock()
	for _, g := range t.groups {
		g.SetFinishMatchTimeSec(t2)
	}
}

func (t *TeamMock) GetStartMatchTimeSec() int64 {
	return t.StartMatchTimeSec
}

func (t *TeamMock) GetFinishMatchTimeSec() int64 {
	t.RLock()
	defer t.RUnlock()
	for _, g := range t.groups {
		return g.GetFinishMatchTimeSec()
	}
	return 0
}

func (t *TeamMock) IsAi() bool {
	t.RLock()
	defer t.RUnlock()
	for _, g := range t.groups {
		for _, p := range g.GetPlayers() {
			if p.IsAi() {
				return true
			}
		}
	}
	return false
}

func (t *TeamMock) GetRank() int {
	return t.rank
}

func (t *TeamMock) SetRank(rank int) {
	t.rank = rank
}

func (t *TeamMock) SortPlayerByRank() []Player {
	t.RLock()
	defer t.RUnlock()
	players := make([]Player, 0, 5)
	for _, g := range t.groups {
		players = append(players, g.GetPlayers()...)
	}
	sort.SliceStable(players, func(i, j int) bool {
		return players[i].GetRank() < players[j].GetRank()
	})
	return players
}

func (t *TeamMock) CanFillAi() bool {
	t.RLock()
	defer t.RUnlock()
	for _, g := range t.groups {
		if !g.CanFillAi() {
			return false
		}
	}
	return true
}
