package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

var JobID int

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	wp := NewWorkerPool(ctx, 10)
	wp.Start()

	rl := NewRateLimiter(10)
	rl.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/worker", func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		select {
		case <-ctx.Done():
			http.Error(w, "Server shutdown", http.StatusServiceUnavailable)
			return
		default:
			JobID++
			payload := r.URL.Query().Get("payload")
			if payload == "" {
				payload = "default"
			}

			wp.AddJob(Job{
				ID:      JobID,
				Payload: payload,
			})
			fmt.Println("Job accepted: ", JobID, payload)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
			return
		}
	}()

	//time.Sleep(5 * time.Second)

	<-ctx.Done()
	fmt.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		panic(err)
		return
	}

	wp.Wait()
	fmt.Println("Server shutdown complete")
}
