package repository

import (
	"appointment-service/internal/model"
	"context"
	"time"
)

type AppointmentRepository interface {
	List(ctx context.Context, trainerID int64) ([]model.Appointment, error)
	Create(ctx context.Context, appointment model.Appointment) (*model.Appointment, error)
	Delete(ctx context.Context, id int64) error
	GetTrainerBookings(ctx context.Context, trainerID int64, startsAt, endsAt time.Time) ([]model.Appointment, error)
	GetClientBookings(ctx context.Context, clientID int64, startsAt, endsAt time.Time) ([]model.Appointment, error)
	Close() error
}
