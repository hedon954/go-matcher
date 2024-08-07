package asynq

import (
	"time"

	"github.com/hedon954/go-matcher/pkg/timer"
)

type Timer[T comparable] struct {
}

func NewTimer[T comparable]() *Timer[T] {
	return &Timer[T]{}
}

func (t *Timer[T]) Register(opType timer.OpType, handler func(id T)) {
	// TODO implement me
	panic("implement me")
}

func (t *Timer[T]) Add(opType timer.OpType, id T, delay time.Duration) error {
	// TODO implement me
	panic("implement me")
}

func (t *Timer[T]) Get(opType timer.OpType, id T) *timer.OperationItem[T] {
	// TODO implement me
	panic("implement me")
}

func (t *Timer[T]) GetAll() []*timer.OperationItem[T] {
	// TODO implement me
	panic("implement me")
}

func (t *Timer[T]) Remove(opType timer.OpType, id T) {
	// TODO implement me
	panic("implement me")
}
