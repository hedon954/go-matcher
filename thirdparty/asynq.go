package thirdparty

import (
	"log/slog"
	"runtime"
	"time"

	"github.com/hibiken/asynq"
)

func StartAsynqServer(redisOpt *asynq.RedisClientOpt,
	queues map[string]int, handler map[string]asynq.HandlerFunc, interval time.Duration) {
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

	// mux maps a type to a handler
	mux := asynq.NewServeMux()

	// register handlers...
	for t, h := range handler {
		mux.HandleFunc(t, h)
	}

	// start server
	slog.Info("start asynq server")
	if err := srv.Run(mux); err != nil {
		slog.Error("could not run asynq server", slog.String("err", err.Error()))
	}
}
