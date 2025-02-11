package postgres

import (
	"appointment-service/internal/config"
	"appointment-service/internal/errors"
	"appointment-service/internal/model"
	"context"
	"log/slog"
	"time"
)

// PostgresAppointmentRepository implements AppointmentRepository using Postgres storage
type PostgresAppointmentRepository struct {
	//db *sql.DB
	//logger *slog.Logger
}

// NewPostgresAppointmentRepository creates a new instance of PostgresAppointmentRepository
func New(dbConfig config.DBConfig, logger *slog.Logger) (*PostgresAppointmentRepository, error) {
	// TODO: Implement
	return nil, errors.InternalError("Create method not implemented yet", nil)
}

func (r *PostgresAppointmentRepository) Create(ctx context.Context, appointment model.Appointment) (*model.Appointment, error) {
	// TODO: Implement
	return nil, errors.InternalError("Create method not implemented yet", nil)
}

func (r *PostgresAppointmentRepository) List(ctx context.Context, trainerId int64) ([]model.Appointment, error) {
	// TODO: Implement
	return nil, errors.InternalError("Create method not implemented yet", nil)
}

func (r *PostgresAppointmentRepository) Delete(ctx context.Context, id int64) error {
	// TODO: Implement
	return errors.InternalError("Create method not implemented yet", nil)
}

func (r *PostgresAppointmentRepository) GetTrainerBookings(ctx context.Context, trainerID int64, startsAt, endsAt time.Time) ([]model.Appointment, error) {
	// TODO: Implement
	return nil, errors.InternalError("Create method not implemented yet", nil)
}

func (r *PostgresAppointmentRepository) GetClientBookings(ctx context.Context, trainerID int64, startsAt, endsAt time.Time) ([]model.Appointment, error) {
	// TODO: Implement
	return nil, errors.InternalError("Create method not implemented yet", nil)
}

func (r *PostgresAppointmentRepository) Close() error {
	return errors.InternalError("Create method not implemented yet", nil)
}
