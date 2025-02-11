package service

import (
	"appointment-service/internal/model"
	"context"
	"time"
)

type AppointmentServicer interface {
	List(ctx context.Context, trainerID int64) ([]model.Appointment, error)
	Create(ctx context.Context, appointment model.Appointment) (*model.Appointment, error)
	GetAvailability(ctx context.Context, trainerID int64, windowStartsAt time.Time, windowEndsAt time.Time) ([]model.TimeSlot, error)
}
