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
			if rId%100 == 0 {
				fmt.Printf("| Room[%d] Match successful, cast time %ds, hasAi: %t, team count: %d\n", rId,
					now-tr.GetStartMatchTimeSec(), tr.HasAi(), len(tr.GetTeams()))
			}
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
