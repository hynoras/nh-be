package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"nh-be/internal/config"

	"github.com/gin-gonic/gin"
)

// ListenAndServe starts the HTTP server, registers health check routes, and
// blocks until a SIGINT or SIGTERM signal is received. Upon receiving a
// shutdown signal, it performs a graceful shutdown by:
//
//  1. Setting the shuttingDown flag (so readiness probes return 503)
//  2. Calling srv.Shutdown with a 5-second timeout to drain active connections
//
// This function is synchronous — it blocks until the server has fully stopped.
// Any deferred cleanup in the caller (e.g. service.Close) will execute only
// after this function returns, preserving correct shutdown ordering.
func ListenAndServe(cfg *config.Config, r *gin.Engine, service *Service) {
	shuttingDown := &atomic.Bool{}

	RegisterHealthRoutes(r, HealthDeps{
		SQLDB:        service.SQLDB,
		Redis:        service.Redis,
		RabbitMQ:     service.RabbitMQ,
		ShuttingDown: shuttingDown,
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		slog.Info("Server is running", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "err", err)
		}
	}()

	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server gracefully...")
	shuttingDown.Store(true)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited cleanly")
}
