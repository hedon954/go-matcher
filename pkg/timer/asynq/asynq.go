package asynq

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hedon954/go-matcher/pkg/timer"
	"github.com/hibiken/asynq"
)

const DefaultQueue = "default"

type Timer[T comparable] struct {
	sync.RWMutex
	client       *asynq.Client
	inspector    *asynq.Inspector
	queue        string
	handlers     map[timer.OpType]func(T)
	tasks        map[string]*asynq.TaskInfo
	key2ID       map[string]T
	expireTicker time.Duration
}

type Option[T comparable] func(*Timer[T])

func NewTimer[T comparable](redisOpt *asynq.RedisClientOpt, opts ...Option[T]) *Timer[T] {
	client := asynq.NewClient(redisOpt)
	inspector := asynq.NewInspector(redisOpt)
	t := &Timer[T]{
		client:       client,
		inspector:    inspector,
		queue:        DefaultQueue,
		handlers:     make(map[timer.OpType]func(T)),
		tasks:        make(map[string]*asynq.TaskInfo),
		key2ID:       map[string]T{},
		expireTicker: time.Minute,
	}

	for _, opt := range opts {
		opt(t)
	}

	go func() {
		ticker := time.NewTicker(t.expireTicker)
		for range ticker.C {
			t.GetAll() // delete expired task
		}
	}()

	return t
}

func WithQueueName[T comparable](queue string) Option[T] {
	return func(t *Timer[T]) {
		t.queue = queue
	}
}

func WithExpireTicker[T comparable](interval time.Duration) Option[T] {
	return func(t *Timer[T]) {
		t.expireTicker = interval
	}
}

func (t *Timer[T]) Register(opType timer.OpType, handler func(T)) {
	t.Lock()
	defer t.Unlock()
	t.handlers[opType] = handler
}

func (t *Timer[T]) Add(opType timer.OpType, id T, delay time.Duration) error {
	t.Lock()
	defer t.Unlock()

	handler := t.handlers[opType]
	if handler == nil {
		return fmt.Errorf("unsupported op type: %s", opType)
	}

	// delete old one
	_ = t.remove(opType, id)

	// add new one
	bs, _ := json.Marshal(id)
	taskInfo, err := t.client.Enqueue(
		asynq.NewTask(string(opType), bs),
		asynq.TaskID(taskKey(opType, id)),
		asynq.Queue(t.queue),
		asynq.ProcessIn(delay),
		asynq.Unique(5*time.Minute))
	if err != nil {
		return err
	}

	// save task info
	t.tasks[taskKey(opType, id)] = taskInfo
	t.key2ID[taskKey(opType, id)] = id
	return err
}

func (t *Timer[T]) Get(opType timer.OpType, id T) *timer.OperationItem[T] {
	t.RLock()
	defer t.RUnlock()
	taskInfo, ok := t.tasks[taskKey(opType, id)]
	if !ok {
		return nil
	}
	return &timer.OperationItem[T]{
		OpType:  opType,
		ID:      id,
		RunTime: taskInfo.NextProcessAt,
	}
}

func (t *Timer[T]) GetAll() []*timer.OperationItem[T] {
	t.Lock()
	defer t.Unlock()

	res := make([]*timer.OperationItem[T], 0, len(t.tasks))
	for key, taskInfo := range t.tasks {
		if taskInfo.NextProcessAt.Before(time.Now()) {
			delete(t.tasks, key)
			continue
		}
		res = append(res, &timer.OperationItem[T]{
			OpType:  timer.OpType(taskInfo.Type),
			ID:      t.key2ID[key],
			RunTime: taskInfo.NextProcessAt,
		})
	}
	return res
}

func (t *Timer[T]) Remove(opType timer.OpType, id T) error {
	t.Lock()
	defer t.Unlock()
	return t.remove(opType, id)
}

func (t *Timer[T]) Stop() {
	_ = t.client.Close()
}

func (t *Timer[T]) remove(opType timer.OpType, id T) error {
	if t.handlers[opType] == nil {
		return nil
	}
	return t.inspector.DeleteTask(t.queue, taskKey(opType, id))
}

func taskKey[T comparable](opType timer.OpType, id T) string {
	return fmt.Sprintf("%s-%v", opType, id)
}
