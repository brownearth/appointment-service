package memory

import (
	"appointment-service/internal/errors"
	"appointment-service/internal/model"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// MemoryAppointmentRepository implements AppointmentRepository using in-memory storage
type MemoryAppointmentRepository struct {
	sync.RWMutex
	appointments []model.Appointment
	lastID       int64
	logger       *slog.Logger
}

// NewMemoryAppointmentRepo creates a new instance of InMemoryAppointmentRepo
func New(logger *slog.Logger) *MemoryAppointmentRepository {
	return &MemoryAppointmentRepository{
		appointments: make([]model.Appointment, 0),
		lastID:       0,
		logger:       logger,
	}
}

// Create stores a new appointment and returns the created appointment with its ID
func (r *MemoryAppointmentRepository) Create(ctx context.Context, appointment model.Appointment) (*model.Appointment, error) {
	r.Lock()
	defer r.Unlock()

	// Check context cancellation
	if ctx.Err() != nil {
		return nil, errors.InternalError("context cancelled", ctx.Err())
	}

	// Generate new ID
	r.lastID++
	newAppointment := appointment
	newAppointment.Id = r.lastID

	r.appointments = append(r.appointments, newAppointment)
	return &newAppointment, nil
}

// List retrieves all appointments for a given trainer
func (r *MemoryAppointmentRepository) List(ctx context.Context, trainerId int64) ([]model.Appointment, error) {
	r.RLock()
	defer r.RUnlock()

	// Check context cancellation
	if ctx.Err() != nil {
		return nil, errors.InternalError("context cancelled", ctx.Err())
	}

	var results []model.Appointment
	for _, apt := range r.appointments {
		if apt.TrainerId == trainerId {
			results = append(results, apt)
		}
	}

	return results, nil
}

// Delete removes an appointment
func (r *MemoryAppointmentRepository) Delete(ctx context.Context, id int64) error {
	r.Lock()
	defer r.Unlock()

	if ctx.Err() != nil {
		return errors.InternalError("context cancelled", ctx.Err())
	}

	for i, apt := range r.appointments {
		if apt.Id == id {
			// Remove element by copying all elements after it one position back
			r.appointments = append(r.appointments[:i], r.appointments[i+1:]...)
			return nil
		}
	}

	return errors.NotFoundError(fmt.Sprintf("appointment with ID %d not found", id))
}

func (r *MemoryAppointmentRepository) GetTrainerBookings(ctx context.Context, trainerID int64, startsAt time.Time, endsAt time.Time) ([]model.Appointment, error) {
	var booked []model.Appointment

	for _, apt := range r.appointments {
		if apt.TrainerId == trainerID &&
			!apt.EndTime.Before(startsAt) &&
			!apt.StartTime.After(endsAt) {
			booked = append(booked, apt)
		}
	}

	return booked, nil
}
func (r *MemoryAppointmentRepository) GetClientBookings(ctx context.Context, clientID int64, startsAt time.Time, endsAt time.Time) ([]model.Appointment, error) {
	var booked []model.Appointment

	for _, apt := range r.appointments {
		if apt.UserId == clientID &&
			!apt.EndTime.Before(startsAt) &&
			!apt.StartTime.After(endsAt) {
			booked = append(booked, apt)
		}
	}

	return booked, nil
}

func (r *MemoryAppointmentRepository) Close() error {
	return nil // No-op in memory storage
}
