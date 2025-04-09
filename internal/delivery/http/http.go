package http

import (
	"context"
	"fmt"
	"net/http"
	"shutdown/internal/domain"
	"shutdown/internal/infrastructure"
)

type Handler struct {
	Context     context.Context
	Processor   domain.Processor
	Pool        *infrastructure.WorkerPool
	RateLimiter *infrastructure.RateLimiter
}

func (h *Handler) Work(w http.ResponseWriter, r *http.Request) {
	if !h.RateLimiter.Allow() {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	payload := r.URL.Query().Get("payload")

	job := domain.Job{Payload: payload}

	h.Pool.AddJob(func() {
		_ = h.Processor.Process(job, h.Context)
	})

	fmt.Fprintf(w, "Job submitted, payload: %s", payload)
}
