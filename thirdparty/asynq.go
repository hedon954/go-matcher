package thirdparty

import (
	"runtime"
	"time"

	"github.com/hibiken/asynq"
)

func NewAsynqServer(redisOpt *asynq.RedisClientOpt,
	queues map[string]int, interval time.Duration) *asynq.Server {
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			// Specify how many concurrent workers to use.
			Concurrency: runtime.NumCPU(),
			// Optionally specify multiple queues with different priority.
			Queues: queues,
			// In test cases, set a small delay to reduce the number of Redis calls.
			// NOTE: `interval` at least 1 second.
			DelayedTaskCheckInterval: interval,
			// See the godoc for other configuration options.
		},
	)
	return srv
}
