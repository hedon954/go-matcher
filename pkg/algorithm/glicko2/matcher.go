package glicko2

import (
	"fmt"
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

	NormalQueue *Queue // 普通队列
	TeamQueue   *Queue // 车队专属队列
}

// NewMatcher 是一个匹配器，包含了 TeamQueue 和 NormalQueue 两个匹配队列
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

// AddGroups 添加队伍
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

func (qm *Matcher) Match() {
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case <-qm.quitChan:
			fmt.Println("\n\nGracefully exit...")
			return
		case <-ticker:
			func() {
				defer func() {
					if err := recover(); err != nil {
						qm.errChan <- fmt.Errorf("glicko2 matcher occurs panic: %v", err)
					}
				}()

				// 取出本轮要匹配的队伍
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

				// 判断哪些 group 需要从专属队列从移动到普通队列
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
					}
					if needMove {
						_ = qm.NormalQueue.AddGroups(g)
					} else {
						_ = qm.TeamQueue.AddGroups(g)
					}
				}

				// 将普通队列中上轮没成功匹配的加回去，下轮重新匹配
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

func (qm *Matcher) GetErrChan() chan error {
	return qm.errChan
}

func (qm *Matcher) GetRoomChan() chan Room {
	if qm.TeamQueue != nil {
		return qm.TeamQueue.roomChan
	}
	if qm.NormalQueue != nil {
		return qm.NormalQueue.roomChan
	}
	return nil
}

func (qm *Matcher) SetNowFunc(f func() int64) {
	qm.NormalQueue.nowUnixFunc = f
	qm.TeamQueue.nowUnixFunc = f
}
