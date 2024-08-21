package native

import (
	"fmt"
	"sync"
	"time"

	"github.com/hedon954/go-matcher/pkg/timer"
)

type Timer struct {
	sync.RWMutex
	handlers map[timer.OpType]func(id int64)
	timers   map[string]*time.Timer
	tasks    map[string]*timer.OperationItem[int64]
}

func NewTimer() *Timer {
	t := &Timer{
		handlers: make(map[timer.OpType]func(id int64)),
		timers:   make(map[string]*time.Timer),
		tasks:    make(map[string]*timer.OperationItem[int64]),
	}
	return t
}

func (t *Timer) Start() {}

func (t *Timer) Register(opType timer.OpType, handler func(id int64)) {
	t.Lock()
	defer t.Unlock()
	t.handlers[opType] = handler
}

func (t *Timer) Add(opType timer.OpType, id int64, delay time.Duration) error {
	handler := t.getHandler(opType)
	if handler == nil {
		return fmt.Errorf("unsupported op type: %s", opType)
	}
	tt := time.AfterFunc(delay, func() {
		_ = t.Remove(opType, id)
		handler(id)
	})
	t.saveTimer(opType, id, tt, delay)
	return nil
}

func (t *Timer) Get(opType timer.OpType, id int64) *timer.OperationItem[int64] {
	t.RLock()
	defer t.RUnlock()
	return t.tasks[timerKey(opType, id)]
}

func (t *Timer) GetAll() []*timer.OperationItem[int64] {
	t.RLock()
	defer t.RUnlock()
	res := make([]*timer.OperationItem[int64], 0, len(t.tasks))
	for _, v := range t.tasks {
		res = append(res, v)
	}
	return res
}

func (t *Timer) Remove(opType timer.OpType, id int64) error {
	t.Lock()
	defer t.Unlock()
	tt, ok := t.timers[timerKey(opType, id)]
	if !ok {
		return nil
	}
	tt.Stop()
	delete(t.timers, timerKey(opType, id))
	delete(t.tasks, timerKey(opType, id))
	return nil
}

func (t *Timer) Stop() {}

func (t *Timer) getHandler(opType timer.OpType) func(id int64) {
	t.RLock()
	defer t.RUnlock()
	return t.handlers[opType]
}

func (t *Timer) saveTimer(opType timer.OpType, id int64, tt *time.Timer, delay time.Duration) {
	t.Lock()
	defer t.Unlock()
	old, ok := t.timers[timerKey(opType, id)]
	if ok {
		old.Stop()
	}
	t.timers[timerKey(opType, id)] = tt
	t.tasks[timerKey(opType, id)] = &timer.OperationItem[int64]{
		OpType:  opType,
		ID:      id,
		RunTime: time.Now().Add(delay),
	}
}

func timerKey(opType timer.OpType, id int64) string {
	return fmt.Sprintf("%s-%d", opType, id)
}
