package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"api-testing-kit/server/internal/config"
	"api-testing-kit/server/internal/db"
	"api-testing-kit/server/internal/httpapi"
)

func main() {
	cfg := config.Load()
	startupCtx, startupCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer startupCancel()

	store, err := db.Open(startupCtx, cfg.DatabaseURL, cfg.DatabaseMaxConns)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}
	if store != nil {
		defer store.Close()
		log.Printf("database connection ready (max_conns=%d)", cfg.DatabaseMaxConns)
	} else {
		log.Print("DATABASE_URL not set; starting API without PostgreSQL")
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.NewRouter(httpapi.RouterDeps{Store: store}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("API server listening on http://localhost:%s", cfg.Port)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		_ = srv.Close()
	}
}
