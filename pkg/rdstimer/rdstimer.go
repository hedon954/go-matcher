package rdstimer

import (
	"strconv"
	"sync"
	"time"

	"github.com/hedon954/go-matcher/pkg/graceful"
)

type TimerFunc func(int)

type TimerManager struct {
	sync.RWMutex
	timerMap map[string]*RdsTimer
}

type TimerOperator interface {
	GetFiredTimerInfo(string) map[string]string
	AddTimer(string, int, int)
	DelTimer(string, int)
}

type RdsTimer struct {
	key      string
	operator TimerOperator
	f        TimerFunc
}

func NewRdsTimerManager() *TimerManager {
	m := new(TimerManager)
	m.timerMap = make(map[string]*RdsTimer)
	return m
}

func newRdsTimer(key string, op TimerOperator, f TimerFunc) *RdsTimer {
	t := new(RdsTimer)
	t.key = key
	t.operator = op
	t.f = f
	return t
}

func (rtm *TimerManager) RegisterTimer(key string, op TimerOperator, f TimerFunc, interval int) {
	rtm.Lock()
	defer rtm.Unlock()
	if _, ok := rtm.timerMap[key]; ok {
		return
	}
	t := newRdsTimer(key, op, f)
	rtm.timerMap[key] = t
	graceful.TimeInterval(time.Duration(interval)*time.Millisecond, t.loop)
}

func (t *RdsTimer) loop() {
	infos := t.operator.GetFiredTimerInfo(t.key)
	now := time.Now().Unix()
	for key, val := range infos {
		timestamp, err := strconv.Atoi(val)
		if err != nil {
			continue
		}
		if int64(timestamp) > now {
			continue
		}
		id, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		graceful.Go(func() { t.f(id) })
	}
}

func (rtm *TimerManager) AddTimer(key string, id int, score int) {
	rtm.RLock()
	t, ok := rtm.timerMap[key]
	if !ok {
		rtm.RUnlock()
		return
	}
	rtm.RUnlock()
	t.operator.AddTimer(key, id, score)
}

func (rtm *TimerManager) DelTimer(key string, id int) {
	rtm.RLock()
	t, ok := rtm.timerMap[key]
	if !ok {
		rtm.RUnlock()
		return
	}
	rtm.RUnlock()
	t.operator.DelTimer(key, id)
}
