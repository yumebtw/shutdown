package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

func StartServer(ctx context.Context, handler *Handler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/work", handler.Work)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Server started at :8080")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")

		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	}()

	return srv
}
