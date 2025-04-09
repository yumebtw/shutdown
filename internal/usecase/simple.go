package usecase

import (
	"context"
	"fmt"
	"shutdown/internal/domain"
	"time"
)

type SimpleProcessor struct{}

func NewSimpleProcessor() *SimpleProcessor {
	return &SimpleProcessor{}
}

func (s *SimpleProcessor) Process(job domain.Job, ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled")
	default:
	}

	fmt.Println("processing job with payload", job.Payload)
	time.Sleep(1 * time.Second)
	fmt.Println("job is done, payload:", job.Payload)
	return nil
}
