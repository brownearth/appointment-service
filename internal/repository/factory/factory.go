package factory

import (
	"appointment-service/internal/config"
	"appointment-service/internal/repository"
	"appointment-service/internal/repository/memory"
	"appointment-service/internal/repository/postgres"
	"appointment-service/internal/repository/sqlite3"
	"fmt"
	"log/slog"
)

// NewRepository creates a new repository based on the provided configuration
func NewRepository(cfg *config.Config, logger *slog.Logger) (repository.AppointmentRepository, error) {
	switch cfg.StorageType {
	case config.Postgres:
		repo, err := postgres.New(cfg.DB, logger.With("repository", "postgres"))
		if err != nil {
			return nil, fmt.Errorf("failed to create postgres repository: %w", err)
		}
		return repo, nil

	case config.SqlLite3:
		repo, err := sqlite3.New(cfg.SqlLite3DbFile, logger.With("repository", "sqlite3"))
		if err != nil {
			return nil, fmt.Errorf("failed to create sqlite repository: %w", err)
		}
		return repo, nil

	case config.Memory:
		return memory.New(logger.With("repository", "memory")), nil

	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.StorageType)
	}
}
