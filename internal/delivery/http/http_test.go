package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"shutdown/internal/infrastructure"
	"shutdown/internal/usecase"
	"testing"
)

func TestHandleJob(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processor := usecase.NewSimpleProcessor()
	pool := infrastructure.NewWorkerPool(2, ctx)
	limiter := infrastructure.NewRateLimiter(5)

	handler := Handler{
		Context:     ctx,
		Processor:   processor,
		Pool:        pool,
		RateLimiter: limiter,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/work", handler.Work)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Post(server.URL+"/work", "application/json", nil)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
