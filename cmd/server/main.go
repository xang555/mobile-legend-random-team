package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	docs "github.com/laoitdev/random-ml-team/docs"
	"github.com/laoitdev/random-ml-team/internal/app"
	"github.com/laoitdev/random-ml-team/internal/config"
	handlers "github.com/laoitdev/random-ml-team/internal/http/handlers"
	"github.com/laoitdev/random-ml-team/internal/http/router"
	"github.com/laoitdev/random-ml-team/internal/random"
	"github.com/laoitdev/random-ml-team/pkg/logger"
)

// @title Random ML Team API
// @version 1.0
// @description API for generating random machine learning teams
// @BasePath /api/v1
func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/config.yaml", "path to configuration file")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	logCfg := logger.Config{
		Level:    cfg.Logging.Level,
		Encoding: cfg.Logging.Encoding,
	}

	log, err := logger.New(logCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to configure logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = log.Sync()
	}()

	docs.SwaggerInfo.BasePath = "/api/v1"

	generator, err := random.NewGenerator(cfg.Team.Composition, cfg.Team.AllowDuplicates, cfg.Team.Heroes)
	if err != nil {
		log.Fatal("invalid team configuration", zap.Error(err))
	}

	teamHandler := handlers.NewTeamHandler(generator, log)
	httpRouter := router.New(log, teamHandler)

	application := app.New(cfg, log, httpRouter)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Info("starting HTTP server", zap.String("addr", application.Addr()))
		errCh <- application.Run()
	}()

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			log.Fatal("server error", zap.Error(err))
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	} else {
		log.Info("server stopped gracefully")
	}
}
