package asynq

import (
	"context"
	"encoding/json"
	"log"
	"sync/atomic"
	"time"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

type EmailDeliveryPayload struct {
	UserID     int
	TemplateID string
}

type ImageResizePayload struct {
	SourceURL string
}

// ----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
// ----------------------------------------------

func NewEmailDeliveryTask(userID int, tmplID string) *asynq.Task {
	payload, _ := json.Marshal(EmailDeliveryPayload{UserID: userID, TemplateID: tmplID})
	return asynq.NewTask(TypeEmailDelivery, payload)
}

func NewImageResizeTask(src string) *asynq.Task {
	payload, _ := json.Marshal(ImageResizePayload{SourceURL: src})
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(TypeImageResize, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute))
}

// ---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
// ---------------------------------------------------------------
func (processor *EmailProcessor) HandleEmailDeliveryTask(_ context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	_ = json.Unmarshal(t.Payload(), &p)
	log.Printf("Handle Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	processor.FinishTaskCount.Add(1)
	return nil
}

func (processor *ImageProcessor) HandleImageResizeTask(_ context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	_ = json.Unmarshal(t.Payload(), &p)
	log.Printf("Resizing image: src=%s", p.SourceURL)
	// Image resizing code ...
	processor.FinishTaskCount.Add(1)
	return nil
}

// ImageProcessor implements asynq.Handler interface.
type ImageProcessor struct {
	FinishTaskCount atomic.Int64
}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return processor.HandleImageResizeTask(ctx, t)
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// EmailProcessor implements asynq.Handler interface.
type EmailProcessor struct {
	FinishTaskCount atomic.Int64
}

func (processor *EmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return processor.HandleEmailDeliveryTask(ctx, t)
}

func NewEmailProcessor() *EmailProcessor {
	return &EmailProcessor{}
}
