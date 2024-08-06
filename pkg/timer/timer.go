package timer

import (
	"time"
)

type OpType string

type OperationItem struct {
	OpType  OpType
	ID      string
	Handler func(id string)
}

type Operator interface {
	Register(opType OpType, handler func(id string))
	Add(opType OpType, id string, delay time.Duration) error
	Remove(opType OpType, id string)
}
