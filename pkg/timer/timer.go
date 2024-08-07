// Package timer defines the interface for timer.
package timer

import (
	"time"
)

// OpType specifies one handler.
type OpType string

// OperationItem specifies one task in timer.
type OperationItem[T comparable] struct {
	// OpType specifies the operation type.
	OpType OpType

	// ID specifies the task unique id.
	ID T

	// RunTime specifies the task run time.
	RunTime time.Time
}

// Operator is the interface for timer.
type Operator[T comparable] interface {
	// Register binds a handler for one operation type.
	Register(opType OpType, handler func(id T))

	// Add adds a new task to timer.
	Add(opType OpType, id T, delay time.Duration) error

	// Get gets the task from timer.
	Get(opType OpType, id T) *OperationItem[T]

	// GetAll gets all tasks from timer.
	GetAll() []*OperationItem[T]

	// Remove removes the task from timer.
	Remove(opType OpType, id T) error

	// Stop stops the timer
	Stop()
}
