package main

import (
	"appointment-service/internal/app"
	"appointment-service/internal/config"
	"appointment-service/internal/logger"
	"appointment-service/internal/version"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration, get version info and initialize logger
	// -----------------------------------------------------------
	cfg := config.Load()
	versionInfo := version.GetInfo()
	logger := initLogger(cfg, &versionInfo)

	// Create application container
	// ----------------------------
	app, err := app.New(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer app.Close()

	// Set up simple signal handling
	// -----------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in anon goroutine
	// ------------------------------
	go func() {
		logger.Info("Starting server", slog.String("port", cfg.Port))
		if err := app.Server.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			logger.Error("Server failed", "error", err)
		}
	}()

	// Wait for interrupt signal
	// -------------------------
	<-quit
	logger.Info("Shutting down server...")

	// Give running requests 5 seconds to complete
	// -------------------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	return nil
}

func initLogger(cfg *config.Config, version *version.Info) *slog.Logger {
	logConfig := logger.Config{
		Attributes: logger.Attributes{
			ServiceName:    "appointment-service",
			ServiceVersion: version.Version,
			CommitSha:      version.Commit,
			BuildTime:      version.BuildTime,
		},
		Level:     logger.ParseLogLevel(cfg.LogLevel, slog.LevelInfo),
		AddSource: cfg.LogSource,
		Format:    logger.ParseFormat(cfg.LogFormat, logger.FormatText),
	}
	return logger.NewLogger(logConfig)
}
