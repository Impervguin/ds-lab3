package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Impervguin/ds-lab2/gateway/internal/adapter/httpclient"
	"github.com/Impervguin/ds-lab2/gateway/internal/handler/health"
	handlerlibrary "github.com/Impervguin/ds-lab2/gateway/internal/handler/library"
	handlerrating "github.com/Impervguin/ds-lab2/gateway/internal/handler/rating"
	handlerreservation "github.com/Impervguin/ds-lab2/gateway/internal/handler/reservation"
	"github.com/Impervguin/ds-lab2/gateway/internal/logger"
	"github.com/Impervguin/ds-lab2/gateway/internal/usecase"
)

const shutdownTimeout = 10 * time.Second

func main() {
	if err := logger.Init(logger.FromEnv("gateway")); err != nil {
		fmt.Fprintf(os.Stderr, "configure logging: %v\n", err)
		os.Exit(1)
	}

	log := logger.Named("bootstrap")
	log.Info("gateway service starting")

	if err := run(log); err != nil {
		log.Error("gateway service stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("gateway service stopped")
}

func run(log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:              net.JoinHostPort("", env("HTTP_PORT", "8080")),
		Handler:           newRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverFailed := make(chan error, 1)
	go func() {
		log.Info("listening", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverFailed <- fmt.Errorf("serve http: %w", err)
		}
		close(serverFailed)
	}()

	select {
	case err := <-serverFailed:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}
	return nil
}

func newRouter() chi.Router {
	libraries := httpclient.NewLibraryClient(env("LIBRARY_SERVICE_URL", "http://localhost:8060"))
	reservations := httpclient.NewReservationClient(env("RESERVATION_SERVICE_URL", "http://localhost:8070"))
	ratings := httpclient.NewRatingClient(env("RATING_SERVICE_URL", "http://localhost:8050"))

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	health.NewHealthHandler().Register(router)
	handlerlibrary.NewLibraryHandler(usecase.NewLibraryUseCase(libraries)).Register(router)
	handlerrating.NewRatingHandler(usecase.NewRatingUseCase(ratings)).Register(router)
	handlerreservation.NewReservationHandler(
		usecase.NewReservationUseCase(reservations, libraries, ratings),
	).Register(router)

	return router
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
