package benchmark

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2/example"
)

// TestQueueMatch 测试匹配性能
/**
遍历数据
    queue_test.go:79: test-0-queue, group count: 100, cast: 2ms, remain group: 55
    queue_test.go:79: test-1-queue, group count: 1000, cast: 88ms, remain group: 274
    queue_test.go:79: test-2-queue, group count: 10000, cast: 8029ms, remain group: 2362
    queue_test.go:79: test-3-queue, group count: 100000, cast: 7733917ms, remain group: 273

二分查找（findGroupForTeam）
    queue_test.go:93: test-0-queue, group count: 100, cast: 1ms, remain group: 69
    queue_test.go:93: test-1-queue, group count: 1000, cast: 48ms, remain group: 354
    queue_test.go:93: test-2-queue, group count: 10000, cast: 3477ms, remain group: 3493
    queue_test.go:91: test-3-queue, group count: 100000, cast: 518299ms, remain group: 50214

二分查找（findGroupForTeam, findTeamForRoom）
    queue_test.go:98: test-0-queue, group count: 100, cast: 1ms, remain group: 25
    queue_test.go:98: test-1-queue, group count: 1000, cast: 12ms, remain group: 153
    queue_test.go:98: test-2-queue, group count: 10000, cast: 793ms, remain group: 1427
    queue_test.go:98: test-3-queue, group count: 100000, cast: 174776ms, remain group: 12787

    TODO: 为什么 100000 的时候性能降这么多？
	前期以为太频繁 groups = append(groups[:i], groups[i+1:]...) 会增加 gc，影响性能，
	所以没有每一次都进行删除操作，而是打标记，待一次匹配完成后再一次性删除，
	这就造成了每次都要遍历一整个 groups，且每次都要 GetState() 去抢锁，反而影响了性能。

遍历查找 + Groups直接从数组移除，不用脏标记
    queue_test.go:124: test-0-queue, group count: 100, cast: 1ms, remain group: 22
    queue_test.go:124: test-1-queue, group count: 1000, cast: 64ms, remain group: 61
    queue_test.go:124: test-2-queue, group count: 10000, cast: 7216ms, remain group: 496

遍历查找 + Groups直接从数据移除 + 利用 groups 的有序性，发现 diffMMR 变大，则提前退出（√）
    queue_test.go:124: test-0-queue, group count: 100, cast: 2ms, remain group: 22
    queue_test.go:124: test-1-queue, group count: 1000, cast: 56ms, remain group: 61
    queue_test.go:124: test-2-queue, group count: 10000, cast: 6271ms, remain group: 496
    queue_test.go:124: test-3-queue, group count: 100000, cast: 1110513ms, remain group: 5567

二分查找 + Groups直接从数组移除，不用脏标记（√）
    queue_test.go:118: test-0-queue, group count: 100, cast: 0ms, remain group: 27
    queue_test.go:118: test-1-queue, group count: 1000, cast: 6ms, remain group: 113
    queue_test.go:118: test-2-queue, group count: 10000, cast: 99ms, remain group: 1095
    queue_test.go:118: test-3-queue, group count: 100000, cast: 3528ms, remain group: 11678
    queue_test.go:109: test-4-queue, group count: 1000000, cast: 270288ms, remain group: 117372
*/
func TestQueueMatch(t *testing.T) {
	groupsData := resolveGroupsData(t)

	roomChan := make(chan glicko2.Room, 128)
	errChan := make(chan error, 128)

	go func() {
		for {
			select {
			case room := <-roomChan:
				// teams := room.GetTeams()
				// t.Logf("room[%d] mmr: %.2f", room.GetID(), room.GetMMR())
				// for _, team := range teams {
				// 	t.Logf("      team mmr: %.2f", team.GetMMR())
				// 	for _, group := range team.GetGroups() {
				// 		t.Logf("             group[%s] mmr: %.2f", group.GetID(), group.GetMMR())
				// 	}
				// }
				room = room
			case e := <-errChan:
				t.Error("has err", e)
			}
		}
	}()

	for index, groupData := range groupsData {
		groups := make([]glicko2.Group, len(groupData))
		for j := 0; j < len(groups); j++ {
			groups[j] = groupData[j]
			groups[j].SetState(glicko2.GroupStateQueuing)
			groups[j].SetStartMatchTimeSec(time.Now().Unix())
		}
		q, _ := glicko2.NewQueue(fmt.Sprintf("test-%d-queue", index), roomChan, example.GetQueueArgs, example.NewTeam,
			example.NewRoom, example.NewRoomWithAi, func() int64 {
				return time.Now().Unix()
			})
		st := time.Now().UnixMilli()
		sameTimes := 0
		oldLen := len(groups)

		for sameTimes < 10 {
			groups = q.Match(groups)
			q.Groups = make(map[string]glicko2.Group, len(groups))
			q.AddGroups(groups...)
			groups = q.GetAndClearGroups()
			ts, fts, rs := q.GetTmpTeamAndFullTeamAndTmpRoom()
			newLen := len(groups)
			for _, tt := range ts {
				newLen += len(tt.GetGroups())
			}
			for _, tt := range fts {
				newLen += len(tt.GetGroups())
			}
			for _, rr := range rs {
				for _, tt := range rr.GetTeams() {
					newLen += len(tt.GetGroups())
				}
			}
			if oldLen == newLen {
				sameTimes++
			} else {
				oldLen = newLen
				sameTimes = 0
			}
		}

		et := time.Now().UnixMilli()
		t.Logf("%s, group count: %d, cast: %dms, remain group: %d", q.Name, len(groupData), et-st, oldLen)

		// groups = q.GetAndClearGroups()
		// ts, rs := q.GetTmpTeamAndTmpRoom()
		// for _, team := range ts {
		// 	groups = append(groups, team.GetGroups()...)
		// }
		// for _, r := range rs {
		// 	for _, team := range r.GetTeams() {
		// 		groups = append(groups, team.GetGroups()...)
		// 	}
		// }
		// for _, group := range groups {
		// 	t.Logf("[%s] player count: %d, mmr: %.2f, start match time: %d\n", group.GetID(), len(group.GetPlayers()),
		// 		group.GetMMR(), group.GetStartMatchTimeSec())
		// }
	}
}

// resolveGroupsData 解析 groups 数据
func resolveGroupsData(t *testing.T) [][]*example.Group {
	// groupCounts := []int{100, 1000}
	groupCounts := []int{100, 1000, 10000, 100000, 1000000}
	res := make([][]*example.Group, 0)
	for _, gc := range groupCounts {
		f, err := os.Open(fmt.Sprintf("groups-%d.json", gc))
		if err != nil {
			t.Error(err)
			continue
		}
		bs, err := io.ReadAll(f)
		if err != nil {
			t.Error(err)
			continue
		}
		groups := make([]*example.Group, 0)
		if err = json.Unmarshal(bs, &groups); err != nil {
			t.Error(err)
			continue
		}
		res = append(res, groups)
	}
	return res
}
