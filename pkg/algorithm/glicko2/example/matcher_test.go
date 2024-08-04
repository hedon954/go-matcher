package example

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

func Test_Matcher(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var roomId = atomic.Int64{}

	errChan := make(chan error, 128)
	roomChan := make(chan glicko2.Room, 128)

	qm, _ := glicko2.NewMatcher(errChan, roomChan, GetQueueArgs, NewTeam, NewRoom, NewRoomWithAi)

	// 异步随机生成 group
	// go func() {
	for i := 0; i < 1000; i++ {
		var players []*Player
		count := rand.Intn(2) + 1
		for j := 0; j < count; j++ {
			p := NewPlayer(uuid.NewString(), false, 0, 1+rand.Intn(200),
				glicko2.Args{
					MMR: 0 + float64(rand.Intn(2000)),
					RD:  0,
					V:   0,
				})
			players = append(players, p)
		}
		newGroup := NewGroup(fmt.Sprintf("Group%d", i+1), players)
		qm.AddGroups(newGroup)
		// ssec := rand.Intn(200)
		// time.Sleep(time.Duration(ssec) * time.Millisecond)
	}
	// }()

	// 异步启动匹配
	go qm.Match()

	// 进程退出
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	for {
		select {
		// 模拟消费 room
		case tr := <-roomChan:
			now := time.Now().Unix()
			rId := roomId.Add(1)
			fmt.Println("-------------------------------------------------------------------")
			fmt.Printf("| Room[%d] Match successful, cast time %ds, hasAi: %t, team count: %d\n", rId,
				now-tr.GetStartMatchTimeSec(), tr.HasAi(), len(tr.GetTeams()))
			for j, team := range tr.GetTeams() {
				fmt.Printf("|   Team %d MMR: %.2f, Star: %d, isAi: %t, cost time %ds\n", j+1,
					team.GetMMR(), team.GetStar(), team.IsAi(), now-team.GetStartMatchTimeSec())
				for _, group := range team.GetGroups() {
					group.SetState(glicko2.GroupStateMatched)
					fmt.Printf("|     %s MMR: %.2f, Star: %d, player count: %d, team type: %d, cost time %ds\n",
						group.GetID(),
						group.GetMMR(),
						group.GetStar(),
						len(group.GetPlayers()), group.Type(),
						now-group.GetStartMatchTimeSec())
					for _, player := range group.GetPlayers() {
						fmt.Printf("|         %s MMR: %.2f, Star: %d, cost time %ds\n",
							player.GetID(),
							player.GetMMR(),
							player.GetStar(),
							now-player.GetStartMatchTimeSec())
					}
				}
			}
			fmt.Println("-------------------------------------------------------------------")
			fmt.Println()
		case err := <-errChan:
			fmt.Println("something error: ", err)
		case <-ch:
			gs1, gs2 := qm.Stop()

			sort.Slice(gs1, func(i, j int) bool {
				return gs1[i].GetMMR() < gs1[j].GetMMR()
			})
			sort.Slice(gs2, func(i, j int) bool {
				return gs2[i].GetMMR() < gs2[j].GetMMR()
			})

			fmt.Println()
			fmt.Println()
			fmt.Println("--------------- finish --------------")

			fmt.Println("normal queue left group count:", len(gs1))
			fmt.Printf("\t\tGroupId\t\t\tPlayerCount\t\tMMR\t\tMatchTime\t\t\n")
			for _, g := range gs1 {
				g.(*Group).Print()
			}
			fmt.Println()
			fmt.Println("team queue left group count:", len(gs2))
			fmt.Printf("\t\tGroupId\t\t\tPlayerCount\t\tMMR\t\tMatchTime\t\t\n")
			for _, g := range gs2 {
				g.(*Group).Print()
			}
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
