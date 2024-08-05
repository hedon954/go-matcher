package glicko2

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

const (
	TeamQueue   = "TeamQueue"
	NormalQueue = "NormalQueue"
)

type Matcher struct {
	errChan  chan error
	quitChan chan struct{}

	NormalQueue *Queue
	TeamQueue   *Queue
}

// NewMatcher is a matcher, which contains both normal queue and team queue.
func NewMatcher(
	errChan chan error,
	roomChan chan Room,
	getQueueArgs func() *QueueArgs,
	newTeamFunc func(group Group) Team,
	newRoomFunc func(team Team) Room,
	newRoomWithAiFunc func(team Team) Room,
) (*Matcher, error) {

	nq, err := NewQueue(NormalQueue, roomChan, getQueueArgs, newTeamFunc, newRoomFunc, newRoomWithAiFunc,
		nowUnicFunc)
	if err != nil {
		return nil, err
	}
	tq, err := NewQueue(TeamQueue, roomChan, getQueueArgs, newTeamFunc, newRoomFunc, newRoomWithAiFunc,
		nowUnicFunc)
	if err != nil {
		return nil, err
	}

	return &Matcher{
		errChan:     errChan,
		quitChan:    make(chan struct{}),
		NormalQueue: nq,
		TeamQueue:   tq,
	}, nil
}

func nowUnicFunc() int64 {
	return time.Now().Unix()
}

// AddGroups adds groups to the queue.
func (qm *Matcher) AddGroups(gs ...Group) error {
	for _, g := range gs {
		groupType := g.Type()
		g.SetState(GroupStateQueuing)
		if groupType == GroupTypeNotTeam {
			if err := qm.NormalQueue.AddGroups(g); err != nil {
				return err
			}
		} else {
			if err := qm.TeamQueue.AddGroups(g); err != nil {
				return err
			}
		}
	}
	return nil
}

func (qm *Matcher) Match(interval time.Duration) {
	ticker := time.NewTicker(interval).C
	for {
		select {
		case <-qm.quitChan:
			slog.Info("stop glicko2 matcher")
			return
		case <-ticker:
			slog.Info("glicko2 matcher tick")
			func() {
				defer func() {
					if err := recover(); err != nil {
						qm.errChan <- fmt.Errorf("glicko2 matcher occurs panic: %v", err)
					}
				}()
				// Get the groups to be matched in this round.
				nGs := qm.NormalQueue.GetAndClearGroups()
				tGs := qm.TeamQueue.GetAndClearGroups()

				wg := sync.WaitGroup{}
				wg.Add(2)
				go func() {
					defer wg.Done()
					nGs = qm.NormalQueue.Match(nGs)
				}()
				go func() {
					defer wg.Done()
					tGs = qm.TeamQueue.Match(tGs)
				}()
				wg.Wait()

				// Determine which groups need to be moved from the exclusive team queue to the normal queue
				now := time.Now()
				for _, g := range tGs {
					needMove := false
					matchTime := now.Unix() - g.GetStartMatchTimeSec()
					switch g.Type() {
					case GroupTypeMaliciousTeam:
						if matchTime >= qm.TeamQueue.MaliciousTeamWaitTimeSec {
							needMove = true
						}
					case GroupTypeUnfriendlyTeam:
						if matchTime >= qm.TeamQueue.UnfriendlyTeamWaitTimeSec {
							needMove = true
						}
					case GroupTypeNormalTeam:
						if matchTime >= qm.TeamQueue.NormalTeamWaitTimeSec {
							needMove = true
						}
					default:
						needMove = false
					}
					if needMove {
						_ = qm.NormalQueue.AddGroups(g)
					} else {
						_ = qm.TeamQueue.AddGroups(g)
					}
				}

				// Add the normal groups back to the normal queue
				_ = qm.NormalQueue.AddGroups(nGs...)
			}()
		}
	}
}

func (qm *Matcher) Stop() ([]Group, []Group) {
	gs1 := qm.NormalQueue.StopMatch()
	gs2 := qm.TeamQueue.StopMatch()
	qm.quitChan <- struct{}{}
	return gs1, gs2
}
