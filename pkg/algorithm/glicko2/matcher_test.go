package glicko2

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Matcher(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var roomId = atomic.Int64{}

	errChan := make(chan error, 128)
	roomChan := make(chan Room, 128)

	qm, _ := NewMatcher(errChan, roomChan, GetQueueArgs, NewTeam, NewRoom, NewRoomWithAi)

	for i := 0; i < 1000; i++ {
		var players []*PlayerMock
		count := rand.Intn(5) + 1
		for j := 0; j < count; j++ {
			p := NewPlayer(uuid.NewString(), false, 0,
				Args{
					MMR: 0 + float64(rand.Intn(2000)),
					RD:  0,
					V:   0,
				})
			players = append(players, p)
		}
		newGroup := NewGroup(fmt.Sprintf("Group%d", i+1), players)
		_ = qm.AddGroups(newGroup)
	}

	// start to match asynchronously
	go qm.Match(time.Millisecond * 100)

	ch := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Second)
		ch <- struct{}{}
	}()
	for {
		select {
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
					group.SetState(GroupStateMatched)
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
			assert.Nil(t, err)
		case <-ch:
			_, _ = qm.Stop()
			fmt.Println("--------------- finish --------------")
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
