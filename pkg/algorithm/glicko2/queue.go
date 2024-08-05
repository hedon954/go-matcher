package glicko2

import (
	"errors"
	"math"
	"sort"
	"sync"
)

const (
	// 每 5 轮会刷新一下配置
	refreshTurn = 5

	// 允许的 mmr 差值百分比，当差值在这个百分比内，就不尝试继续寻找更优解了
	acceptMMRDiffPercent = 0.01

	// 使用二分查找的阈值，在寻找匹配的 group 和 team 的时候，有两种思路：
	// 1. 遍历查找：找到最匹配的那个（数据量大的时候性能很低）
	// 2. 二分查找：二分查找加局部搜索寻找局部最优解（性能很高，但可能错误最优解）
	// 所以在数量量小的时候，使用遍历寻找最优解，数据量大的时候，使用二分寻找局部最优解
	useBinarySearchThreshold = 1000

	CancelMatchByServerStop = "Failed to match. Please try again later"
	CancelMatchByTimeout    = "No team found, please try again later"
)

var (
	ErrQueueClosed      = errors.New("Match queue has been closed")
	ErrNilGetArgsFunc   = errors.New("the func getQueueArgs is nil")
	ErrNilGetArgsReturn = errors.New("getQueueArgs() == nil")
)

// Queue 是一个匹配队列
type Queue struct {
	lock          sync.Mutex
	isClosed      bool                   // 是否已关闭
	Name          string                 // 队列名称
	Groups        map[string]Group       // 在队列中的队伍，对于 Groups 的所有处理都要加锁
	FullTeam      []Team                 // 匹配过程的满员临时阵营，只在 Match 中调用，不可以并发调用
	TmpTeam       []Team                 // 匹配过程中的临时阵营，只能在 Match 中调用，不可以并发调用
	TmpRoom       []Room                 // 匹配过程中的临时房间，只能在 Match 中调用，不可以并发调用
	roomChan      chan Room              // 匹配成功的房间会投进这个 channel
	newTeam       func(group Group) Team // 构建新 team 的方法
	newRoom       func(team Team) Room   // 构建新 room 的方法
	newRoomWithAi func(team Team) Room   // 构建带 ai 的新 room 的方法
	nowUnixFunc   func() int64           // 返回当前时间的时间戳
	matchTurn     int                    // 匹配轮次，对 5 取模，用于定时刷新配置
	*QueueArgs                           // 队列参数
	getQueueArgs  func() *QueueArgs      // 队列参数获取方法，用于定时刷新配置
}

type QueueArgs struct {
	MatchTimeoutSec int64 `json:"match_timeout_sec"` // 匹配超时时间

	TeamPlayerLimit int `json:"team_player_limit"` // 阵营总人数上限
	RoomTeamLimit   int `json:"room_team_limit"`   // 房间总阵营数
	MinPlayerCount  int `json:"min_player_count"`  // TODO: 最小开局人数

	NewerWithNewer bool `json:"newer_with_newer"` // 新手是否只和新手匹配到一起

	UnfriendlyTeamMMRVarianceMin int `json:"unfriendly_team_mmr_variance_min"` // 不友好车队 mmr 最小方差
	MaliciousTeamMMRVarianceMin  int `json:"malicious_team_mmr_variance_min"`  // 恶意车队 mmr 最小方差

	NormalTeamWaitTimeSec     int64 `json:"normal_team_wait_time_sec"`     // 普通车队在专属队列中的匹配时长
	UnfriendlyTeamWaitTimeSec int64 `json:"unfriendly_team_wait_time_sec"` // 不友好车队在专属队列中的匹配时长
	MaliciousTeamWaitTimeSec  int64 `json:"malicious_team_wait_time_sec"`  // 恶意车队在专属队列中的匹配时长

	MatchRanges []MatchRange `json:"match_ranges"` // 匹配范围策略
}

type MatchRange struct {
	MaxMatchSec   int64 `json:"max_match_sec"`   // 最长匹配时间 s（不包含）
	MMRGapPercent int   `json:"mmr_gap_percent"` // 允许的 mmr 差距百分比(0~100)（包含），0 表示无限制
	CanJoinTeam   bool  `json:"can_join_team"`   // 是否加入满人车队
	StarGap       int   `json:"star_gap"`        // 允许的段位差距数（包含），0 表示无限制
}

var defaultMatchRange = MatchRange{
	MaxMatchSec:   15,
	MMRGapPercent: 10,
	CanJoinTeam:   false,
	StarGap:       12,
}

func NewQueue(
	name string, roomChan chan Room,
	getQueueArgs func() *QueueArgs,
	newTeamFunc func(group Group) Team,
	newRoomFunc func(team Team) Room,
	newRoomWithAiFunc func(team Team) Room,
	nowUnixFunc func() int64,
) (*Queue, error) {
	if getQueueArgs() == nil {
		return nil, ErrNilGetArgsFunc
	}
	args := getQueueArgs()
	if args == nil {
		return nil, ErrNilGetArgsReturn
	}
	return &Queue{
		lock:          sync.Mutex{},
		Name:          name,
		roomChan:      roomChan,
		Groups:        make(map[string]Group, 128),
		TmpTeam:       make([]Team, 0, 128),
		FullTeam:      make([]Team, 0, 128),
		TmpRoom:       make([]Room, 0, 128),
		newTeam:       newTeamFunc,
		newRoom:       newRoomFunc,
		newRoomWithAi: newRoomWithAiFunc,
		getQueueArgs:  getQueueArgs,
		QueueArgs:     args,
		nowUnixFunc:   nowUnixFunc,
	}, nil
}

func (q *Queue) SortedGroups() []Group {
	groups := q.AllGroups()
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].GetMMR() < groups[j].GetMMR()
	})
	return groups
}

func (q *Queue) AllGroups() []Group {
	q.Lock()
	defer q.Unlock()

	groups := make([]Group, len(q.Groups))
	index := 0
	for _, group := range q.Groups {
		groups[index] = group
		index++
	}
	return groups
}

// AddGroups 添加 group 到队列中
func (q *Queue) AddGroups(gs ...Group) error {
	q.Lock()
	defer q.Unlock()

	if q.isClosed {
		return ErrQueueClosed
	}

	for _, g := range gs {
		if g.GetStartMatchTimeSec() == 0 {
			g.SetStartMatchTimeSec(q.nowUnixFunc())
		}
		q.Groups[g.GetID()] = g
	}
	return nil
}

// GetAndClearGroups 取出要匹配的 group 并且清空当前 groups 列表
func (q *Queue) GetAndClearGroups() []Group {
	q.Lock()
	defer q.Unlock()
	now := q.nowUnixFunc()
	res := make([]Group, 0, len(q.Groups))
	for _, g := range q.Groups {
		// 只要还在匹配中的队伍
		if g.GetState() == GroupStateQueuing {
			// 去掉超时的队伍
			if q.MatchTimeoutSec != 0 && now-g.GetStartMatchTimeSec() >= q.MatchTimeoutSec {
				waitSec := q.nowUnixFunc() - g.GetStartMatchTimeSec()
				g.SetState(GroupStateUnready)
				g.SetStartMatchTimeSec(0)
				tmpG := g
				go func() {
					tmpG.ForceCancelMatch(CancelMatchByTimeout, waitSec)
				}()
				continue
			}
			res = append(res, g)
		}
	}
	q.Groups = make(map[string]Group, 128)
	return res
}

// clearTmp 清除临时数据并归位 groups
func (q *Queue) clearTmp() []Group {
	groups := make([]Group, 0, 128)
	for _, t := range q.TmpTeam {
		for _, g := range t.GetGroups() {
			if g.GetState() == GroupStateQueuing {
				groups = append(groups, g)
			} else {
				g.SetStartMatchTimeSec(0)
			}
		}
	}
	for _, t := range q.FullTeam {
		for _, g := range t.GetGroups() {
			if g.GetState() == GroupStateQueuing {
				groups = append(groups, g)
			} else {
				g.SetStartMatchTimeSec(0)
			}
		}
	}
	for _, r := range q.TmpRoom {
		for _, rt := range r.GetTeams() {
			for _, g := range rt.GetGroups() {
				if g.GetState() == GroupStateQueuing {
					groups = append(groups, g)
				} else {
					g.SetStartMatchTimeSec(0)
				}
			}
		}
	}
	q.TmpTeam = q.TmpTeam[:0]
	q.FullTeam = q.FullTeam[:0]
	q.TmpRoom = q.TmpRoom[:0]
	return groups
}

// Match 队列匹配逻辑
func (q *Queue) Match(groups []Group) []Group {
	q.Lock()
	defer q.Unlock()
	// 对 groups 进行排序，后面用二分查找提高效率
	sortGroupsByMMR(groups)
	// 构建新的 team
	groups = q.buildNewTeams(groups)
	// 优先填充 AI
	q.fillTeamsWithAi()
	// 对 team 排序
	q.sortFullTeamsByMMR()
	// 创建新的房间
	q.buildNewRooms()
	// 每轮都打散重来，每 refreshTurn 轮刷新 QueueArgs
	q.refreshMatchTurn()
	// 清除临时数据，这里不保留临时数据是怕后面 n 轮匹配的时候，临时数据中的 Group 已经取消匹配了
	gs := q.clearTmp()
	groups = append(groups, gs...)
	return groups
}

func (q *Queue) buildNewTeams(groups []Group) []Group {
	// 构建新的 team
	tryTeamCounts := len(groups)
	var found bool
	for i := 0; len(groups) > 0 && i < tryTeamCounts; i++ {
		var team Team
		for gPos, g := range groups {
			if g.GetState() == GroupStateQueuing {
				team = q.newTeam(g)
				groups = append(groups[:gPos], groups[gPos+1:]...)
				break
			}
		}
		if team == nil {
			break
		}

		for team.PlayerCount() < q.TeamPlayerLimit {
			groups, found = q.findGroupForTeam(team, groups)
			if !found {
				break
			}
		}
		if team.IsFull(q.TeamPlayerLimit) {
			q.FullTeam = append(q.FullTeam, team)
		} else {
			q.TmpTeam = append(q.TmpTeam, team)
		}
	}
	return groups
}

func (q *Queue) buildNewRooms() {
	tryRoomTimes := len(q.FullTeam)
	for l := 0; len(q.FullTeam) > 0 && l < tryRoomTimes; l++ {
		room := q.newRoom(q.FullTeam[0])
		q.FullTeam = q.FullTeam[1:]
		for len(room.GetTeams()) < q.RoomTeamLimit {
			if !q.findTeamForRoom(room) {
				break
			}
		}
		if len(room.GetTeams()) >= q.RoomTeamLimit {
			q.roomMatchSuccess(room)
			continue
		}
		q.TmpRoom = append(q.TmpRoom, room)
	}
}

// fillTeamsWithAi 填充 AI
func (q *Queue) fillTeamsWithAi() {
	newFullTeam := make([]Team, 0, len(q.FullTeam))
	for _, team := range q.FullTeam {
		if team.CanFillAi() {
			newRoom := q.newRoomWithAi(team)
			q.roomMatchSuccess(newRoom)
			continue
		}
		newFullTeam = append(newFullTeam, team)
	}
	q.FullTeam = newFullTeam
}

func (q *Queue) refreshMatchTurn() {
	q.matchTurn = (q.matchTurn + 1) % refreshTurn
	if q.matchTurn == 0 && q.getQueueArgs != nil {
		newQueueArgs := q.getQueueArgs()
		if newQueueArgs != nil {
			q.QueueArgs = newQueueArgs
		}
	}
}

// sortGroupsByMMR 根据 MMR 对 groups 进行排序
func sortGroupsByMMR(groups []Group) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].GetMMR() < groups[j].GetMMR()
	})
}

// sortFullTeamsByMMR 根据 MMR 对 TmpTeam 进行排序
func (q *Queue) sortFullTeamsByMMR() {
	sort.Slice(q.FullTeam, func(i, j int) bool {
		return q.FullTeam[i].GetMMR() < q.FullTeam[j].GetMMR()
	})
}

// findGroupForTeam 从 groups 中找到适合 team 的 group 并加入其中
// 1. 符合条件
// 2. mmr 尽可能接近
func (q *Queue) findGroupForTeam(team Team, groups []Group) ([]Group, bool) {
	if len(groups) == 0 {
		return groups, false
	}
	if len(groups) > useBinarySearchThreshold {
		return q.findGroupForTeamBinary(team, groups)
	}
	return q.findGroupForTeamRange(team, groups)
}

// findGroupForTeamRange 从 groups 中找到适合 team 的 group 并加入其中（遍历）
func (q *Queue) findGroupForTeamRange(team Team, groups []Group) ([]Group, bool) {
	// 第1个队伍直接进
	if team.PlayerCount() == 0 {
		for gPos, group := range groups {
			if group.GetState() == GroupStateQueuing {
				team.AddGroup(group)
				groups = append(groups[:gPos], groups[gPos+1:]...)
				return groups, true
			}
		}
		return groups, false
	}

	// 寻找平均 mmr 最接近的 group 组成一个 team
	teamMMR := team.GetMMR()
	teamPlayerCount := team.PlayerCount()
	closestIndex := -1
	minMMRDiff := math.MaxFloat64
	for i, group := range groups {
		if group.GetState() != GroupStateQueuing {
			continue
		}
		if teamPlayerCount+group.PlayerCount() > q.TeamPlayerLimit {
			continue
		}
		mmrDiff := math.Abs(group.GetMMR() - teamMMR)
		if mmrDiff < minMMRDiff {
			if !q.canGroupTogether(team, group) {
				continue
			}
			minMMRDiff = mmrDiff
			closestIndex = i
		} else if mmrDiff >= minMMRDiff && closestIndex != -1 {
			// 因为 groups 是按 mmr 升序的，
			// 所以如果 mmrDiff 在变大的话，说明后面的 mmrDiff 只会更大，
			// 所以这里如果已经找到 group 了，那么是可以提前结束的
			break
		}
	}

	if closestIndex == -1 {
		return groups, false
	}

	team.AddGroup(groups[closestIndex])
	groups = append(groups[:closestIndex], groups[closestIndex+1:]...)
	return groups, true
}

// findGroupForTeamBinary 从 groups 中找到适合 team 的 group 并加入其中（二分查找）
func (q *Queue) findGroupForTeamBinary(team Team, groups []Group) ([]Group, bool) {
	// 第1个队伍直接进
	if team.PlayerCount() == 0 {
		for gPos, group := range groups {
			if group.GetState() == GroupStateQueuing {
				team.AddGroup(group)
				groups = append(groups[:gPos], groups[gPos+1:]...)
				return groups, true
			}
		}
		return groups, false
	}

	// 寻找平均 mmr 最接近的 group 组成一个 team
	teamMMR := team.GetMMR()
	teamPlayerCount := team.PlayerCount()
	closestIndex := -1

	// 二分查找定位一个接近点
	low, high := 0, len(groups)-1
	minDifference := math.MaxFloat64
	for low <= high {
		mid := low + (high-low)/2
		group := groups[mid]
		groupMMR := group.GetMMR()
		mmrDifference := math.Abs(groupMMR - teamMMR)
		if mmrDifference < minDifference && teamPlayerCount+group.PlayerCount() <= q.TeamPlayerLimit && group.GetState() == GroupStateQueuing {
			if q.canGroupTogether(team, group) {
				minDifference = mmrDifference
				closestIndex = mid
			}
		}
		if groupMMR < teamMMR {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	// 局部线性搜索来确认是否有更接近的匹配
	searchRange := 10 // 根据实际情况调整这个范围
	start := int(math.Max(0, float64(closestIndex-searchRange)))
	end := int(math.Min(float64(len(groups)-1), float64(closestIndex+searchRange)))
	for i := start; i <= end; i++ {
		group := groups[i]
		if teamPlayerCount+group.PlayerCount() > q.TeamPlayerLimit || group.GetState() != GroupStateQueuing {
			continue
		}

		mmrDifference := math.Abs(group.GetMMR() - teamMMR)
		if mmrDifference < minDifference && q.canGroupTogether(team, group) {
			minDifference = mmrDifference
			closestIndex = i
		}
		if closestIndex != -1 && mmrDifference > minDifference {
			break
		}
	}

	if closestIndex == -1 {
		return groups, false
	}

	team.AddGroup(groups[closestIndex])
	groups = append(groups[:closestIndex], groups[closestIndex+1:]...)
	return groups, true
}

// findTeamForRoom 从 TmpTeam 中找到合适 room 的 team 并加入其中
func (q *Queue) findTeamForRoom(room Room) bool {
	if len(q.FullTeam) == 0 {
		return false
	}
	if len(q.FullTeam) > useBinarySearchThreshold {
		return q.findTeamForRoomBinary(room)
	}
	return q.findTeamForRoomRange(room)
}

// findTeamForRoomRange 从 TmpTeam 中找到合适 room 的 team 并加入其中（遍历）
func (q *Queue) findTeamForRoomRange(room Room) bool {
	if len(q.FullTeam) == 0 {
		return false
	}

	teams := room.GetTeams()

	// 已经满了，直接返回
	if len(teams) >= q.RoomTeamLimit {
		return false
	}

	// 如果 room 是空的，则第一个 team 直接进入 room 中
	if len(teams) == 0 {
		room.AddTeam(q.FullTeam[0])
		q.FullTeam = q.FullTeam[1:]
		return true
	}

	// 寻找 mmr 最接近的 team
	roomMMR := room.GetMMR()
	closestIndex := -1
	minMMRDiff := math.MaxFloat64
	for i, team := range q.FullTeam {
		// 只有当 team 已经组建完毕了，才可以加入到 room 中
		if team.PlayerCount() != q.TeamPlayerLimit {
			continue
		}
		mmrDiff := math.Abs(team.GetMMR() - roomMMR)
		if mmrDiff < minMMRDiff {
			if !q.canTeamTogether(room, team) {
				continue
			}
			minMMRDiff = mmrDiff
			closestIndex = i
		} else if mmrDiff > minMMRDiff && closestIndex != -1 {
			break
		}
	}

	if closestIndex == -1 {
		return false
	}

	room.AddTeam(q.FullTeam[closestIndex])
	q.FullTeam = append(q.FullTeam[:closestIndex], q.FullTeam[closestIndex+1:]...)
	return true
}

// findTeamForRoomBinary 从 TmpTeam 中找到合适 room 的 team 并加入其中（二分查找）
func (q *Queue) findTeamForRoomBinary(room Room) bool {
	if len(q.FullTeam) == 0 {
		return false
	}

	teams := room.GetTeams()

	// 已经满了，直接返回
	if len(teams) >= q.RoomTeamLimit {
		return false
	}

	// 如果 room 是空的，则第一个 team 直接进入 room 中
	if len(teams) == 0 {
		room.AddTeam(q.FullTeam[0])
		q.FullTeam = q.FullTeam[1:]
		return true
	}

	// 二分查找定位一个接近点，寻找 mmr 最接近的 team
	roomMMR := room.GetMMR()
	closestIndex := -1
	low, high := 0, len(q.FullTeam)-1
	minDifference := math.MaxFloat64
	for low <= high {
		mid := low + (high-low)/2
		team := q.FullTeam[mid]
		teamMMR := team.GetMMR()
		mmrDifference := math.Abs(roomMMR - teamMMR)
		if mmrDifference < minDifference && q.canTeamTogether(room, team) {
			minDifference = mmrDifference
			closestIndex = mid
			// 如果只相差1%，那就算可以了，不继续寻找了
			if minDifference <= roomMMR*acceptMMRDiffPercent {
				break
			}
		}
		if teamMMR < roomMMR {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	// 只相差1%，那就算可以了，不继续寻找了
	if closestIndex != -1 && minDifference <= roomMMR*acceptMMRDiffPercent {
		room.AddTeam(q.FullTeam[closestIndex])
		q.FullTeam = append(q.FullTeam[:closestIndex], q.FullTeam[closestIndex+1:]...)
		return true
	}

	// 局部线性搜索来确认是否有更接近的匹配
	searchRange := 10 // 根据实际情况调整这个范围
	start := int(math.Max(0, float64(closestIndex-searchRange)))
	end := int(math.Min(float64(len(q.FullTeam)-1), float64(closestIndex+searchRange)))

	for i := start; i <= end; i++ {
		team := q.FullTeam[i]
		mmrDifference := math.Abs(team.GetMMR() - roomMMR)
		if mmrDifference < minDifference && q.canTeamTogether(room, team) {
			minDifference = mmrDifference
			closestIndex = i
		}
	}

	if closestIndex == -1 {
		return false
	}

	room.AddTeam(q.FullTeam[closestIndex])
	q.FullTeam = append(q.FullTeam[:closestIndex], q.FullTeam[closestIndex+1:]...)
	return true
}

// canGroupTogether 判断队伍之间是否可以组成一个阵营
func (q *Queue) canGroupTogether(team Team, group Group) bool {
	ngMMR := group.GetMMR()
	ngStar := group.GetStar()
	ngCanFillAi := group.CanFillAi()
	ngIsNewer := group.IsNewer()

	for _, g := range team.GetGroups() {
		// 优先把可以匹配到 AI 的匹配到一起
		if g.CanFillAi() != ngCanFillAi {
			return false
		}

		if q.NewerWithNewer && (ngIsNewer != g.IsNewer()) {
			return false
		}

		// 获取匹配范围配置
		mr := q.getMatchRange(g.GetStartMatchTimeSec(), group.GetStartMatchTimeSec())

		// mmr 是否匹配
		gMMR := g.GetMMR()
		if mr.MMRGapPercent != 0 && math.Abs(gMMR-ngMMR) > gMMR*float64(mr.MMRGapPercent)/100 {
			return false
		}

		// 段位是否匹配
		if mr.StarGap != 0 && int(math.Abs(float64(g.GetStar()-ngStar))) > mr.StarGap {
			return false
		}
	}
	return true
}

// canTeamTogether 判断阵营之间是否可以组成一个房间
func (q *Queue) canTeamTogether(room Room, tt Team) bool {
	ttMMR := tt.GetMMR()
	ttStar := tt.GetStar()
	ttNewer := tt.IsNewer()
	// 判断 tt 是否满足跟当前 room 中的所有 team 匹配的条件
	// 只要有一个不满足，就返回 false
	for _, t := range room.GetTeams() {
		mr := q.getMatchRange(t.GetStartMatchTimeSec(), tt.GetStartMatchTimeSec())

		if q.NewerWithNewer && (ttNewer != t.IsNewer()) {
			return false
		}

		// 是否加入车队
		if len(t.GetGroups()) > 1 && !mr.CanJoinTeam && len(tt.GetGroups()) == 1 {
			return false
		}

		// mmr 是否匹配
		tMMR := t.GetMMR()
		if mr.MMRGapPercent != 0 && math.Abs(tMMR-ttMMR) > tMMR*float64(mr.MMRGapPercent)/float64(100) {
			return false
		}

		// 段位是否匹配
		if mr.StarGap != 0 && int(math.Abs(float64(t.GetStar()-ttStar))) > mr.StarGap {
			return false
		}
	}
	return true
}

// getMatchRange 获取匹配范围
func (q *Queue) getMatchRange(mst1, mst2 int64) MatchRange {
	now := q.nowUnixFunc()

	if len(q.MatchRanges) == 0 {
		return defaultMatchRange
	}

	// 以匹配时间短的那个为准
	mt := int64(math.Min(float64(now-mst1), float64(now-mst2)))
	for _, mr := range q.MatchRanges {
		if mt < mr.MaxMatchSec {
			return mr
		}
	}

	// 默认返回最后一个
	return q.MatchRanges[len(q.MatchRanges)-1]
}

// StopMatch 取消匹配
func (q *Queue) StopMatch() []Group {
	q.Lock()
	defer q.Unlock()
	q.isClosed = true

	groups := q.clearTmp()
	for _, group := range groups {
		q.Groups[group.GetID()] = group
	}
	remainGroups := make([]Group, 0, len(groups))
	for _, g := range q.Groups {
		if g.GetState() != GroupStateQueuing {
			continue
		}
		waitSec := q.nowUnixFunc() - g.GetStartMatchTimeSec()
		g.SetState(GroupStateUnready)
		g.SetStartMatchTimeSec(0)
		g.ForceCancelMatch(CancelMatchByServerStop, waitSec)
		remainGroups = append(remainGroups, g)
	}
	q.Groups = make(map[string]Group)
	return remainGroups
}

// roomMatchSuccess 房间匹配成功
func (q *Queue) roomMatchSuccess(room Room) {
	go func() {
		q.setRoomReady(room)
		q.roomChan <- room
	}()
}

// setRoomReady 更新 group 状态为匹配完成
func (q *Queue) setRoomReady(room Room) {
	for _, t := range room.GetTeams() {
		for _, g := range t.GetGroups() {
			g.SetState(GroupStateMatched)
		}
	}
}

func (q *Queue) Lock() {
	q.lock.Lock()
}

func (q *Queue) Unlock() {
	q.lock.Unlock()
}
