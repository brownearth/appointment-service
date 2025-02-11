package sqlite3

import (
	"appointment-service/internal/errors"
	"appointment-service/internal/model"
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// New creates a new SQLite3 appointment repository with the given database path.
// Returns error if connection fails.
func New(dbPath string, logger *slog.Logger) (*Repository, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	log.Printf("Connected to SQLite DB at: %s", dbPath)
	return &Repository{db: db, logger: logger}, nil
}

// Create inserts a new appointment into the database.
// Returns the created appointment with generated ID or error if insert fails.
func (r *Repository) Create(ctx context.Context, apt model.Appointment) (*model.Appointment, error) {
	const query = `
		INSERT INTO appointments (trainer_id, user_id, start_time, end_time)
		VALUES (:trainer_id, :user_id, :start_time, :end_time)
		RETURNING id, trainer_id, user_id, start_time, end_time`

	log.Printf("Creating appointment: %+v", apt)

	// Convert domain model to DB model
	dbApt := toDBModel(apt)
	rows, err := r.db.NamedQueryContext(ctx, query, dbApt)
	if err != nil {
		return nil, fmt.Errorf("creating appointment: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned after insert")
	}

	// Scan the returned row into DB model
	var created dbAppointment
	if err := rows.StructScan(&created); err != nil {
		return nil, fmt.Errorf("scanning created appointment: %w", err)
	}

	result := toDomainModel(created)
	log.Printf("Created appointment: %+v", result)
	return &result, nil
}

// List retrieves all appointments for a given trainer ID.
// Returns empty slice if no appointments found.
func (r *Repository) List(ctx context.Context, trainerID int64) ([]model.Appointment, error) {
	const query = `
		SELECT id, trainer_id, user_id, start_time, end_time
		FROM appointments
		WHERE trainer_id = ?`

	var dbAppts []dbAppointment
	if err := r.db.SelectContext(ctx, &dbAppts, query, trainerID); err != nil {
		return nil, fmt.Errorf("listing appointments: %w", err)
	}

	appointments := toDomainModels(dbAppts)
	log.Printf("Listed %d appointments for trainer %d", len(appointments), trainerID)
	return appointments, nil
}

// Delete removes an appointment by ID.
// Returns NotFoundError if appointment doesn't exist.
func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM appointments WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting appointment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking affected rows: %w", err)
	}

	if rows == 0 {
		return errors.NotFoundError(fmt.Sprintf("appointment %d not found", id))
	}

	log.Printf("Deleted appointment %d", id)
	return nil
}

// GetTrainerBookings retrieves all appointments for a trainer within the given time range.
// Time range is inclusive of start and end times.
func (r *Repository) GetTrainerBookings(ctx context.Context, trainerID int64, start, end time.Time) ([]model.Appointment, error) {
	const query = `
		SELECT id, trainer_id, user_id, start_time, end_time
		FROM appointments
		WHERE trainer_id = ?
		AND end_time >= ?
		AND start_time <= ?`

	var dbAppts []dbAppointment
	if err := r.db.SelectContext(ctx, &dbAppts, query, trainerID, start, end); err != nil {
		return nil, fmt.Errorf("getting booked appointments: %w", err)
	}

	appointments := toDomainModels(dbAppts)
	return appointments, nil
}

// GetClientBookings retrieves all appointments for a user within the given time range.
// Time range is inclusive of start and end times.
func (r *Repository) GetClientBookings(ctx context.Context, userId int64, start, end time.Time) ([]model.Appointment, error) {
	const query = `
		SELECT id, trainer_id, user_id, start_time, end_time
		FROM appointments
		WHERE user_id = ?
		AND end_time >= ?
		AND start_time <= ?`

	var dbAppts []dbAppointment
	if err := r.db.SelectContext(ctx, &dbAppts, query, userId, start, end); err != nil {
		return nil, fmt.Errorf("getting booked appointments: %w", err)
	}

	appointments := toDomainModels(dbAppts)
	return appointments, nil
}

// Close closes the database connection.
func (r *Repository) Close() error {
	return r.db.Close()
}
