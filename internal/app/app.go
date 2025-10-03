package app

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/laoitdev/random-ml-team/internal/config"
)

// App represents the HTTP server instance.
type App struct {
	cfg    *config.Config
	server *http.Server
	logger *zap.Logger
}

// New creates an application instance ready to be started.
func New(cfg *config.Config, logger *zap.Logger, handler http.Handler) *App {
	return &App{
		cfg: cfg,
		logger: logger,
		server: &http.Server{
			Addr:         cfg.Addr(),
			Handler:      handler,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

// Run starts the HTTP server and blocks until it exits.
func (a *App) Run() error {
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the HTTP server.
func (a *App) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, a.cfg.Server.ShutdownTimeout)
	defer cancel()

	return a.server.Shutdown(ctx)
}

// Addr exposes the address the server listens on.
func (a *App) Addr() string {
	return a.server.Addr
}

// WaitForShutdown wraps the graceful shutdown logic with logging.
func (a *App) WaitForShutdown(ctx context.Context, cancel context.CancelFunc, errCh <-chan error) error {
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// GracefulStop orchestrates a shutdown with timeout.
func (a *App) GracefulStop(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, a.cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("graceful shutdown failed", zap.Error(err))
	}
	// Give outstanding connections time to close.
	time.Sleep(200 * time.Millisecond)
}
