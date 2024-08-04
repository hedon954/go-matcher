package mocktest

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2/example"

	glicko "github.com/zelenin/go-glicko2"
)

/**
算法效果验证：
	设定模拟1000个机器人，每个机器人设定获得不同名次的概率权重
	机器人匹配到一起时，队伍的最终获得名次的权重为该队伍的平均权重
	使用新的匹配机制，按照mmr初始匹配，胜负规则判定按照队伍的平均名次概率
	判定名次后变化分数，连续跑1000轮，导出如下数据：
	1. 机器人名称
	2. 机器人设定胜率
	3. 机器人实际胜率
	4. 机器人当前mmr分数
	5. 机器人当前的段位
*/

// Test_runMatchAndSettle 执行匹配和结算并输出结果
func Test_runMatchAndSettle(t *testing.T) {
	players := resolveRobots()
	groups := playersToGroups(players)
	roomChan := make(chan glicko2.Room, 128)

	go func() {
		var roomId = atomic.Int64{}
		for {
			select {
			case tr := <-roomChan:
				fmt.Println()
				now := time.Now().Unix()
				rId := roomId.Add(1)
				fmt.Println("-------------------------------------------------------------------")
				fmt.Printf("| Room[%d] Match successful, cast time %ds, hasAi: %t, team count: %d\n", rId,
					now-tr.GetStartMatchTimeSec(), tr.HasAi(), len(tr.GetTeams()))
				for j, team := range tr.GetTeams() {
					fmt.Printf("|   Team %d average MMR: %.2f, isAi: %t, cost time %ds\n", j+1,
						team.GetMMR(), team.IsAi(), now-team.GetStartMatchTimeSec())
					for _, group := range team.GetGroups() {
						group.SetState(glicko2.GroupStateMatched)
						fmt.Printf("|     %s MMR: %.2f, player count: %d, team type: %d, cost time %ds\n",
							group.GetID(),
							group.GetMMR(),
							len(group.GetPlayers()), group.Type(),
							now-group.GetStartMatchTimeSec())
					}
				}
				fmt.Println("-------------------------------------------------------------------")
				calcuteRoomResult(tr.(*example.Room))
				fmt.Println()
			}
		}
	}()

	// 比赛轮次
	battleCount := 1000
	for i := 0; i < battleCount; i++ {
		q, _ := glicko2.NewQueue("RobotQueue", roomChan, example.GetQueueArgs, example.NewTeam, example.NewRoom,
			example.NewRoomWithAi, func() int64 { return time.Now().Unix() })
		q.AddGroups(groups...)
		gs := q.GetAndClearGroups()
		q.Match(gs)
		if i%100 == 0 {
			fname := "mock-robots-result-" + strconv.Itoa(i) + ".json"
			os.Remove(fname)
			f, err := os.Create(fname)
			if err != nil {
				t.Error(err)
				continue
			}
			sort.Slice(players, func(i, j int) bool {
				return players[i].GetMMR() > players[j].GetMMR()
			})
			bs, _ := json.Marshal(players)
			f.Write(bs)
			f.Close()
		}
	}

	time.Sleep(5 * time.Second)
	// 保存结果
}

// Test_genreateRobots 生成机器人数据
func Test_genreateRobots(t *testing.T) {
	robotCount := 1000
	players := make([]*example.Player, robotCount)
	for i := 0; i < robotCount; i++ {
		p := example.NewPlayer(fmt.Sprintf("Player-%d", i+1), false, 0, 0, glicko2.Args{
			MMR: 1500,
			RD:  360,
			V:   0.06,
		})
		p.WinWeight = 10 + rand.Intn(90)
		players[i] = p
	}
	bs, err := json.Marshal(players)
	if err != nil {
		t.Fatal(err)
		return
	}
	f, err := os.Create("mock-robots.json")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer f.Close()
	_, _ = f.Write(bs)
}

// resolveRobots 解析机器人数据
func resolveRobots() []*example.Player {
	var res []*example.Player
	f, err := os.Open("mock-robots.json")
	if err != nil {
		panic(err)
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(bs, &res); err != nil {
		panic(err)
	}
	for _, p := range res {
		p.Player = glicko.NewPlayer(glicko.NewRating(p.MMR, p.RD, p.V))
	}
	return res
}

// playersToGroups 每个 player 一个 Group
func playersToGroups(players []*example.Player) []glicko2.Group {
	res := make([]glicko2.Group, len(players))
	for i := 0; i < len(players); i++ {
		g := example.NewGroup("Group-"+strconv.Itoa(i+1), []*example.Player{players[i]})
		g.SetState(glicko2.GroupStateQueuing)
		res[i] = g
	}
	return res
}

// calcuteRoomResult 计算房间的胜负，根据 WinWeight
func calcuteRoomResult(room *example.Room) {
	// 算出每个 Team 的权重
	for _, team := range room.GetTeams() {
		t := team.(*example.Team)
		for _, g := range team.GetGroups() {
			players := g.GetPlayers()
			for _, p := range players {
				player := p.(*example.Player)
				t.WinWeight += player.WinWeight
			}
		}
	}
	// 阵营间随机胜负
	ts := room.GetTeams()
	teams := make([]*example.Team, len(ts))
	for i, t := range ts {
		teams[i] = t.(*example.Team)
	}
	rank := 1
	var t *example.Team
	for len(teams) > 0 {
		t, teams = weightedRandomSelectTeam(teams)
		t.SetRank(rank)
		rank++
	}
	// 阵营内随机排名
	for _, team := range ts {
		r := 1
		players := make([]*example.Player, 0)
		groups := team.GetGroups()
		for _, g := range groups {
			ps := g.GetPlayers()
			for _, p := range ps {
				players = append(players, p.(*example.Player))
			}
		}
		var p *example.Player
		for len(players) > 0 {
			p, players = weightedRandomSelectPlayer(players)
			p.SetRank(r)
			r++
		}
	}
}

// weightedRandomSelectTeam 随机 team 胜负
func weightedRandomSelectTeam(teams []*example.Team) (*example.Team, []*example.Team) {
	totalWeight := 0
	for _, team := range teams {
		totalWeight += team.WinWeight
	}
	randWeight := rand.Intn(totalWeight)
	for i, team := range teams {
		randWeight -= team.WinWeight
		if randWeight < 0 {
			// 移除选中的团队
			return team, append(teams[:i], teams[i+1:]...)
		}
	}
	panic("should never reach here")
}

// weightedRandomSelectTeam 随机 team 胜负
func Test_weightedRandomSelectTeam(t *testing.T) {
	teams := make([]*example.Team, 3)
	for i := 0; i < len(teams); i++ {
		teams[i] = new(example.Team)
		teams[i].Id = i
	}
	teams[0].WinWeight = 60
	teams[1].WinWeight = 30
	teams[2].WinWeight = 10

	res := make([][]int, 3)
	for i := 0; i < len(res); i++ {
		res[i] = make([]int, 3)
	}
	for i := 0; i < 1000; i++ {
		tmp := make([]*example.Team, 3)
		copy(tmp, teams)
		var team *example.Team
		rank := 0
		for len(tmp) > 0 {
			team, tmp = weightedRandomSelectTeam(tmp)
			res[team.Id][rank]++
			rank++
		}
	}

	t.Log(res)
}

// weightedRandomSelectPlayer 随机 player 排名
func weightedRandomSelectPlayer(players []*example.Player) (*example.Player, []*example.Player) {
	totalWeight := 0
	for _, player := range players {
		totalWeight += player.WinWeight
	}
	randWeight := rand.Intn(totalWeight)
	for i, player := range players {
		randWeight -= player.WinWeight
		if randWeight < 0 {
			// 移除选中的玩家
			return player, append(players[:i], players[i+1:]...)
		}
	}
	panic("should never reach here")
}
