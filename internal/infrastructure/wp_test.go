package infrastructure

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWorkerPool_AddJobAndWait(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	limit := 10

	wp := NewWorkerPool(limit, ctx)

	calls := make(chan int, limit)
	for i := 1; i <= limit; i++ {
		wp.AddJob(func() {
			func(ctx context.Context) {
				select {
				case <-ctx.Done():
					return
				default:
				}
				calls <- 1
			}(ctx)
		})
	}

	select {
	case <-ctx.Done():
		fmt.Println("timeout waiting for jobs to finish")
		return
	default:
	}

	wp.Wait()

	if len(calls) < 10 {
		t.Errorf("expected 10 tasks to run, got %d", len(calls))
	}
}
