// Package timer defines the interface for timer.
package timer

import (
	"time"
)

// OpType specifies one handler.
type OpType string

// OperationItem specifies one task in timer.
type OperationItem struct {
	// OpType specifies the operation type.
	OpType OpType

	// ID specifies the task unique id.
	ID string

	// RunTime specifies the task run time.
	RunTime time.Time
}

// Operator is the interface for timer.
type Operator interface {
	// Register binds a handler for one operation type.
	Register(opType OpType, handler func(id string))

	// Add adds a new task to timer.
	Add(opType OpType, id string, delay time.Duration) error

	// Get gets the task from timer.
	Get(opType OpType, id string) *OperationItem

	// GetAll gets all tasks from timer.
	GetAll() []*OperationItem

	// Remove removes the task from timer.
	Remove(opType OpType, id string)
}
