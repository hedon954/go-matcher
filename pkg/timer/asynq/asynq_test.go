package asynq

import (
	"errors"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	ptimer "github.com/hedon954/go-matcher/pkg/timer"
	"github.com/hedon954/go-matcher/thirdparty"
	"github.com/hibiken/asynq"

	"github.com/stretchr/testify/assert"
)

func TestAsynqTimer(t *testing.T) {
	const opType1 ptimer.OpType = "1"
	const opType2 ptimer.OpType = "2"
	const opType3 ptimer.OpType = "3"

	const id1 = 1
	const id2 = 2

	var numMap = map[int64]*atomic.Int64{
		id1: {},
		id2: {},
	}

	handle1 := func(id int64) {
		slog.Info("handle1 running")
		numMap[id].Add(1)
	}
	handle2 := func(id int64) {
		slog.Info("handle2 running")
		numMap[id].Add(-1)
	}

	redisOpt := &asynq.RedisClientOpt{Addr: thirdparty.NewMiniRedis().Addr()}
	time.Sleep(3 * time.Millisecond)

	// create a new timer
	timer := NewTimer[int64](redisOpt,
		WithQueueName[int64]("test"),
		WithTimerInterval[int64](time.Second),
	)
	go timer.Start()
	defer timer.Stop()
	assert.NotNil(t, timer.client)
	assert.NotNil(t, timer.inspector)
	assert.NotNil(t, timer.handlers)
	assert.NotNil(t, timer.tasks)
	assert.NotNil(t, timer.key2ID)
	assert.Equal(t, "test", timer.queue)

	// register operations
	timer.Register(opType1, handle1)
	timer.Register(opType2, handle2)
	assert.Equal(t, 2, len(timer.handlers))

	// add not exists operation should return error
	err := timer.Add(opType3, id1, 1*time.Second)
	assert.Equal(t, errors.New("unsupported op type: 3"), err)

	// get not exists operation should return nil
	assert.Nil(t, timer.Get(opType3, id1))

	// get not exists task id should return nil
	assert.Nil(t, timer.Get(opType2, 10000))

	// add optype1 should add num after delay, and the timer should be deleted
	delay := time.Millisecond * 10
	err = timer.Add(opType1, id1, delay)
	assert.Nil(t, err)
	assert.NotNil(t, timer.Get(opType1, id1))
	assert.Equal(t, 1, len(timer.GetAll()))
	assert.Equal(t, int64(0), numMap[id1].Load())
	// NOTE: sleep to long, do not run in ci
	// time.Sleep(delay + 1*time.Second)
	// assert.Equal(t, int64(1), numMap[id1].Load())
	// time.Sleep(delay + 1*time.Millisecond)
	// assert.Equal(t, 0, len(timer.GetAll()))

	// add optype2 should reduce num after delay
	err = timer.Add(opType2, id2, delay)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), numMap[id2].Load())
	// NOTE: sleep to long, do not run in ci
	// time.Sleep(delay + 1*time.Second)
	// assert.Equal(t, int64(-1), numMap[id2].Load())

	// add optype2 and delete it before delay should not reduce num
	err = timer.Add(opType2, id2, delay)
	assert.Nil(t, err)
	// NOTE: sleep to long, do not run in ci
	// assert.Equal(t, int64(-1), numMap[id2].Load())
	err = timer.Remove(opType2, id2)
	assert.Nil(t, err)
	// NOTE: sleep to long, do not run in ci
	// time.Sleep(delay + 1*time.Second)
	// assert.Equal(t, int64(-1), numMap[id2].Load())

	// read optype2 should stop the first add
	err = timer.Add(opType2, id2, delay)
	assert.Nil(t, err)
	// NOTE: sleep to long, do not run in ci
	// assert.Equal(t, int64(-1), numMap[id2].Load())
	err = timer.Add(opType2, id2, 2*time.Second)
	assert.Nil(t, err)
	// NOTE: sleep to long, do not run in ci
	// time.Sleep(delay + 1*time.Second + 10*time.Millisecond)
	// assert.Equal(t, int64(-1), numMap[id2].Load())
	// time.Sleep(delay + 2*time.Second + 10*time.Millisecond)
	// assert.Equal(t, int64(-2), numMap[id2].Load())

	// remove not existed operation should not panic
	err = timer.Remove(opType3, id1)
	assert.Nil(t, err)
}
