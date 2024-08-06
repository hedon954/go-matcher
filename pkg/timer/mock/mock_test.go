package mock

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	ptimer "github.com/hedon954/go-matcher/pkg/timer"

	"github.com/stretchr/testify/assert"
)

func TestMockTimer(t *testing.T) {
	const opType1 ptimer.OpType = "1"
	const opType2 ptimer.OpType = "2"
	const opType3 ptimer.OpType = "3"

	const id1 = "id-1"
	const id2 = "id-2"

	var numMap = map[string]*atomic.Int64{
		id1: {},
		id2: {},
	}

	// create a new timer
	timer := NewTimer()
	assert.NotNil(t, timer.timers)
	assert.NotNil(t, timer.handlers)

	// resgiter operations
	timer.Register(opType1, func(id string) {
		numMap[id].Add(1)
	})
	timer.Register(opType2, func(id string) {
		numMap[id].Add(-1)
	})
	assert.Equal(t, 2, len(timer.handlers))

	// add not exists operation should return error
	err := timer.Add(opType3, id1, 1*time.Second)
	assert.Equal(t, errors.New("unsupported op type: 3"), err)

	// add optype1 should add num after delay, and the timer should be deleted
	delay := time.Millisecond * 10
	err = timer.Add(opType1, id1, delay)
	assert.Nil(t, err)
	assert.NotNil(t, timer.Get(opType1, id1))
	assert.Equal(t, 1, len(timer.GetAll()))
	assert.Equal(t, int64(0), numMap[id1].Load())
	time.Sleep(delay + 3*time.Millisecond)
	assert.Equal(t, int64(1), numMap[id1].Load())
	assert.Equal(t, 0, len(timer.timers))

	// add optype2 should reduce num after delay
	err = timer.Add(opType2, id2, delay)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), numMap[id2].Load())
	time.Sleep(delay + 3*time.Millisecond)
	assert.Equal(t, int64(-1), numMap[id2].Load())

	// add optype2 and delete it before delay should not reduce num
	err = timer.Add(opType2, id2, delay)
	assert.Nil(t, err)
	assert.Equal(t, int64(-1), numMap[id2].Load())
	timer.Remove(opType2, id2)
	time.Sleep(delay + 3*time.Millisecond)
	assert.Equal(t, int64(-1), numMap[id2].Load())

	// readd optype2 should stop the first add
	err = timer.Add(opType2, id2, delay)
	assert.Nil(t, err)
	assert.Equal(t, int64(-1), numMap[id2].Load())
	err = timer.Add(opType2, id2, delay+5*time.Millisecond)
	assert.Nil(t, err)
	time.Sleep(delay + 3*time.Millisecond)
	assert.Equal(t, int64(-1), numMap[id2].Load())
	time.Sleep(delay + 3*time.Millisecond + 5*time.Millisecond)
	assert.Equal(t, int64(-2), numMap[id2].Load())

	// remove not existed operation should not panic
	timer.Remove(opType3, id1)
}
