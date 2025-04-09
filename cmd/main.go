package main

import (
	"context"
	"log"
	//"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpDelivery "shutdown/internal/delivery/http"
	"shutdown/internal/infrastructure"
	"shutdown/internal/usecase"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	processor := usecase.NewSimpleProcessor()
	pool := infrastructure.NewWorkerPool(5, ctx)
	limiter := infrastructure.NewRateLimiter(10)

	handler := httpDelivery.Handler{
		Context:     ctx,
		Processor:   processor,
		Pool:        pool,
		RateLimiter: limiter,
	}

	httpServer := httpDelivery.StartServer(ctx, &handler)

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	pool.Wait()
	log.Println("Shutdown complete")
}
