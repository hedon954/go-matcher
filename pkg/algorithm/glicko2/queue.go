package glicko2

import (
	"errors"
	"math"
	"sort"
	"sync"
)

const (
	// Configuration refreshes every 5 turns
	refreshTurn = 5

	// Allowed MMR difference percentage. When the difference is within this percentage,
	// it will not attempt to find a better solution.
	acceptMMRDiffPercent = 0.01

	// Threshold for using binary search. When looking for matching groups and teams,
	// there are two approaches:
	// 1. Sequential search: Find the best match (low performance with large data).
	// 2. Binary search: Use binary search and local search to find the local optimal solution
	//    (high performance, but might miss the optimal solution).
	// Therefore, when the quantity is small, use sequential search to find the optimal solution,
	// and when the quantity is large, use binary search to find the local optimal solution.
	useBinarySearchThreshold = 1000

	CancelMatchByServerStop = "Failed to match. Please try again later"
	CancelMatchByTimeout    = "No team found, please try again later"
)

var (
	ErrQueueClosed      = errors.New("Match queue has been closed")
	ErrNilGetArgsFunc   = errors.New("the func getQueueArgs is nil")
	ErrNilGetArgsReturn = errors.New("getQueueArgs() == nil")
)

// Queue is a match queue
type Queue struct {
	lock          sync.Mutex
	isClosed      bool                   // Whether the queue is closed
	Name          string                 // Queue name
	Groups        map[string]Group       // Groups in the queue, all operations on Groups need to be locked
	FullTeam      []Team                 // Fully formed temporary teams during matching, called only in Match, not thread-safe
	TmpTeam       []Team                 // Temporary teams during matching, called only in Match, not thread-safe
	TmpRoom       []Room                 // Temporary rooms during matching, called only in Match, not thread-safe
	roomChan      chan Room              // Successfully matched rooms are sent to this channel
	newTeam       func(group Group) Team // Method to create a new team
	newRoom       func(team Team) Room   // Method to create a new room
	newRoomWithAi func(team Team) Room   // Method to create a new room with AI
	nowUnixFunc   func() int64           // Function to return the current timestamp
	matchTurn     int                    // Match turn, modulus 5, used to periodically refresh the configuration
	*QueueArgs                           // Queue parameters
	getQueueArgs  func() *QueueArgs      // Method to get queue parameters, used to periodically refresh the configuration
}

type QueueArgs struct {
	MatchTimeoutSec int64 `json:"match_timeout_sec"` // Match timeout duration

	TeamPlayerLimit int `json:"team_player_limit"` // Team player limit
	RoomTeamLimit   int `json:"room_team_limit"`   // Room team limit
	MinPlayerCount  int `json:"min_player_count"`  // TODO: Minimum number of players to start

	NewerWithNewer bool `json:"newer_with_newer"` // Whether beginners can only match with beginners

	UnfriendlyTeamMMRVarianceMin int `json:"unfriendly_team_mmr_variance_min"` // Minimum MMR variance for unfriendly teams
	MaliciousTeamMMRVarianceMin  int `json:"malicious_team_mmr_variance_min"`  // Minimum MMR variance for malicious teams

	NormalTeamWaitTimeSec     int64 `json:"normal_team_wait_time_sec"`     // Matching duration for normal teams in exclusive queue
	UnfriendlyTeamWaitTimeSec int64 `json:"unfriendly_team_wait_time_sec"` // Matching duration for unfriendly teams in exclusive queue
	MaliciousTeamWaitTimeSec  int64 `json:"malicious_team_wait_time_sec"`  // Matching duration for malicious teams in exclusive queue

	MatchRanges []MatchRange `json:"match_ranges"` // Match range strategies
}

type MatchRange struct {
	MaxMatchSec   int64 `json:"max_match_sec"`   // Maximum match duration in seconds (exclusive)
	MMRGapPercent int   `json:"mmr_gap_percent"` // Allowed MMR difference percentage (0-100) (inclusive), 0 means no restriction
	CanJoinTeam   bool  `json:"can_join_team"`   // Whether to join full teams
	StarGap       int   `json:"star_gap"`        // Allowed rank difference (inclusive), 0 means no restriction
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

// AddGroups adds groups to the queue
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

// GetAndClearGroups retrieves and clears the current groups list
func (q *Queue) GetAndClearGroups() []Group {
	q.Lock()
	defer q.Unlock()
	now := q.nowUnixFunc()
	res := make([]Group, 0, len(q.Groups))
	for _, g := range q.Groups {
		// Only groups that are still in the queue
		if g.GetState() == GroupStateQueuing {
			// Remove groups that have timed out
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

// clearTmp clears temporary data and resets groups
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

// Match queue matching logic
func (q *Queue) Match(groups []Group) []Group {
	q.Lock()
	defer q.Unlock()
	// Sort groups for later binary search to improve efficiency
	sortGroupsByMMR(groups)
	// Build new teams
	groups = q.buildNewTeams(groups)
	// Prioritize filling AI
	q.fillTeamsWithAi()
	// Sort teams by MMR
	q.sortFullTeamsByMMR()
	// Create new rooms
	q.buildNewRooms()
	// Shuffle and reset every turn, refresh QueueArgs every refreshTurn
	q.refreshMatchTurn()
	// Clear temporary data, not keeping temporary data is to prevent groups from canceling matches in later turns
	gs := q.clearTmp()
	groups = append(groups, gs...)
	return groups
}

func (q *Queue) buildNewTeams(groups []Group) []Group {
	// Build new teams
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

// fillTeamsWithAi fills AI
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

// sortGroupsByMMR sorts groups by MMR
func sortGroupsByMMR(groups []Group) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].GetMMR() < groups[j].GetMMR()
	})
}

// sortFullTeamsByMMR sorts TmpTeam by MMR
func (q *Queue) sortFullTeamsByMMR() {
	sort.Slice(q.FullTeam, func(i, j int) bool {
		return q.FullTeam[i].GetMMR() < q.FullTeam[j].GetMMR()
	})
}

// findGroupForTeam finds a suitable group for the team from groups and adds it
// 1. Meet the conditions
// 2. MMR is as close as possible
func (q *Queue) findGroupForTeam(team Team, groups []Group) ([]Group, bool) {
	if len(groups) == 0 {
		return groups, false
	}
	if len(groups) > useBinarySearchThreshold {
		return q.findGroupForTeamBinary(team, groups)
	}
	return q.findGroupForTeamRange(team, groups)
}

// findGroupForTeamRange finds a suitable group for the team from groups and adds it (sequential search)
func (q *Queue) findGroupForTeamRange(team Team, groups []Group) ([]Group, bool) {
	// The first group directly joins
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

	// Find the group with the closest average MMR to form a team
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
			// Because groups are sorted in ascending order of MMR,
			// if mmrDiff is increasing, it means that the subsequent mmrDiff will only get larger,
			// so if a group has already been found, it can be terminated early.
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

// findGroupForTeamBinary finds a suitable group for the team from groups and adds it (binary search)
func (q *Queue) findGroupForTeamBinary(team Team, groups []Group) ([]Group, bool) {
	// The first group directly joins
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

	// Find the group with the closest average MMR to form a team
	teamMMR := team.GetMMR()
	teamPlayerCount := team.PlayerCount()
	closestIndex := -1

	// Binary search to locate a close point
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

	// Local linear search to confirm if there is a closer match
	searchRange := 10 // Adjust this range based on actual conditions
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

// findTeamForRoom finds a suitable team for the room from TmpTeam and adds it
func (q *Queue) findTeamForRoom(room Room) bool {
	if len(q.FullTeam) == 0 {
		return false
	}
	if len(q.FullTeam) > useBinarySearchThreshold {
		return q.findTeamForRoomBinary(room)
	}
	return q.findTeamForRoomRange(room)
}

// findTeamForRoomRange finds a suitable team for the room from TmpTeam and adds it (sequential search)
func (q *Queue) findTeamForRoomRange(room Room) bool {
	if len(q.FullTeam) == 0 {
		return false
	}

	teams := room.GetTeams()

	// Already full, return directly
	if len(teams) >= q.RoomTeamLimit {
		return false
	}

	// If the room is empty, the first team directly enters the room
	if len(teams) == 0 {
		room.AddTeam(q.FullTeam[0])
		q.FullTeam = q.FullTeam[1:]
		return true
	}

	// Find the team with the closest MMR
	roomMMR := room.GetMMR()
	closestIndex := -1
	minMMRDiff := math.MaxFloat64
	for i, team := range q.FullTeam {
		// Only fully formed teams can join the room
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

// findTeamForRoomBinary finds a suitable team for the room from TmpTeam and adds it (binary search)
func (q *Queue) findTeamForRoomBinary(room Room) bool {
	if len(q.FullTeam) == 0 {
		return false
	}

	teams := room.GetTeams()

	// Already full, return directly
	if len(teams) >= q.RoomTeamLimit {
		return false
	}

	// If the room is empty, the first team directly enters the room
	if len(teams) == 0 {
		room.AddTeam(q.FullTeam[0])
		q.FullTeam = q.FullTeam[1:]
		return true
	}

	// Binary search to locate a close point, find the team with the closest MMR
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
			// If the difference is only 1%, consider it acceptable and stop searching
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

	// If the difference is only 1%, consider it acceptable and stop searching
	if closestIndex != -1 && minDifference <= roomMMR*acceptMMRDiffPercent {
		room.AddTeam(q.FullTeam[closestIndex])
		q.FullTeam = append(q.FullTeam[:closestIndex], q.FullTeam[closestIndex+1:]...)
		return true
	}

	// Local linear search to confirm if there is a closer match
	searchRange := 10 // Adjust this range based on actual conditions
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

// canGroupTogether determines whether groups can form a team
func (q *Queue) canGroupTogether(team Team, group Group) bool {
	ngMMR := group.GetMMR()
	ngStar := group.GetStar()
	ngCanFillAi := group.CanFillAi()
	ngIsNewer := group.IsNewer()

	for _, g := range team.GetGroups() {
		// Prioritize matching groups that can fill AI together
		if g.CanFillAi() != ngCanFillAi {
			return false
		}

		if q.NewerWithNewer && (ngIsNewer != g.IsNewer()) {
			return false
		}

		// Get match range configuration
		mr := q.getMatchRange(g.GetStartMatchTimeSec(), group.GetStartMatchTimeSec())

		// Check if MMR matches
		gMMR := g.GetMMR()
		if mr.MMRGapPercent != 0 && math.Abs(gMMR-ngMMR) > gMMR*float64(mr.MMRGapPercent)/100 {
			return false
		}

		// Check if rank matches
		if mr.StarGap != 0 && int(math.Abs(float64(g.GetStar()-ngStar))) > mr.StarGap {
			return false
		}
	}
	return true
}

// canTeamTogether determines whether teams can form a room
func (q *Queue) canTeamTogether(room Room, tt Team) bool {
	ttMMR := tt.GetMMR()
	ttStar := tt.GetStar()
	ttNewer := tt.IsNewer()
	// Check if tt meets the matching conditions with all teams in the current room
	// If any condition is not met, return false
	for _, t := range room.GetTeams() {
		mr := q.getMatchRange(t.GetStartMatchTimeSec(), tt.GetStartMatchTimeSec())

		if q.NewerWithNewer && (ttNewer != t.IsNewer()) {
			return false
		}

		// Check if joining full teams is allowed
		if len(t.GetGroups()) > 1 && !mr.CanJoinTeam && len(tt.GetGroups()) == 1 {
			return false
		}

		// Check if MMR matches
		tMMR := t.GetMMR()
		if mr.MMRGapPercent != 0 && math.Abs(tMMR-ttMMR) > tMMR*float64(mr.MMRGapPercent)/float64(100) {
			return false
		}

		// Check if rank matches
		if mr.StarGap != 0 && int(math.Abs(float64(t.GetStar()-ttStar))) > mr.StarGap {
			return false
		}
	}
	return true
}

// getMatchRange gets the match range
func (q *Queue) getMatchRange(mst1, mst2 int64) MatchRange {
	now := q.nowUnixFunc()

	if len(q.MatchRanges) == 0 {
		return defaultMatchRange
	}

	// Use the shorter match duration as the standard
	mt := int64(math.Min(float64(now-mst1), float64(now-mst2)))
	for _, mr := range q.MatchRanges {
		if mt < mr.MaxMatchSec {
			return mr
		}
	}

	// Default return the last one
	return q.MatchRanges[len(q.MatchRanges)-1]
}

// StopMatch cancels the match
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

// roomMatchSuccess indicates a successful room match
func (q *Queue) roomMatchSuccess(room Room) {
	go func() {
		q.setRoomReady(room)
		q.roomChan <- room
	}()
}

// setRoomReady updates the group status to match completed
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
