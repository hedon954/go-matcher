package mock

import (
	"fmt"
	"sync"
	"time"

	"github.com/hedon954/go-matcher/pkg/timer"
)

type Timer struct {
	sync.RWMutex
	handlers map[timer.OpType]func(id string)
	timers   map[string]*time.Timer
}

func NewTimer() *Timer {
	t := &Timer{
		handlers: make(map[timer.OpType]func(id string)),
		timers:   make(map[string]*time.Timer),
	}
	return t
}

func (t *Timer) Register(opType timer.OpType, handler func(id string)) {
	t.Lock()
	defer t.Unlock()
	t.handlers[opType] = handler
}

func (t *Timer) Add(opType timer.OpType, id string, delay time.Duration) error {
	handler := t.getHandler(opType)
	if handler == nil {
		return fmt.Errorf("unsupported op type: %s", opType)
	}
	tt := time.AfterFunc(delay, func() {
		handler(id)
	})
	t.saveTimer(opType, id, tt)
	return nil
}

func (t *Timer) Remove(opType timer.OpType, id string) {
	t.Lock()
	defer t.Unlock()
	tt, ok := t.timers[timerKey(opType, id)]
	if !ok {
		return
	}
	tt.Stop()
	delete(t.timers, timerKey(opType, id))
}

func (t *Timer) getHandler(opType timer.OpType) func(id string) {
	t.RLock()
	defer t.RUnlock()
	return t.handlers[opType]
}

func (t *Timer) saveTimer(opType timer.OpType, id string, tt *time.Timer) {
	t.Lock()
	defer t.Unlock()
	old, ok := t.timers[timerKey(opType, id)]
	if ok {
		old.Stop()
	}
	t.timers[timerKey(opType, id)] = tt
}

func timerKey(opType timer.OpType, id string) string {
	return fmt.Sprintf("%s-%s", opType, id)
}
