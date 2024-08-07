package asynq

import (
	"log"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/hibiken/asynq"
)

func TestAsync_shouldwork(t *testing.T) {
	redis := newMiniRedis()
	// redis := "127.0.0.1:6379"
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redis})

	ep := NewEmailProcessor()
	ip := NewImageProcessor()
	go startHandleServer(ep, ip, redis)

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	task, err := NewEmailDeliveryTask(42, "some:template:id")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------

	info, err = client.Enqueue(task, asynq.ProcessIn(3*time.Millisecond))
	if err != nil {
		log.Fatalf("could not schedule task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	task, err = NewImageResizeTask("https://example.com/myassets/image.jpg")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task, asynq.ProcessIn(3*time.Millisecond), asynq.MaxRetry(10),
		asynq.Timeout(10*time.Millisecond))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	time.Sleep(10 * time.Second)
}

func startHandleServer(ep *EmailProcessor, ip *ImageProcessor, redis string) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redis},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 2,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// In test cases, set a small delay to reduce the number of Redis calls
			DelayedTaskCheckInterval: 1 * time.Millisecond,
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, ep.HandleEmailDeliveryTask)
	mux.Handle(TypeImageResize, ip)
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func newMiniRedis() string {
	var miniRedisClient *miniredis.Miniredis
	var err error
	miniRedisClient, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	return miniRedisClient.Addr()
}
