package service

import (
	"appointment-service/internal/errors"
	"appointment-service/internal/model"
	"appointment-service/internal/repository"
	"context"
	"fmt"
	"log/slog"
	"time"
)

const appointmentDuration = 30 * time.Minute

type AppointmentService struct {
	repo   repository.AppointmentRepository
	logger *slog.Logger
}

func NewAppointmentService(repo repository.AppointmentRepository, logger *slog.Logger) AppointmentServicer {
	return &AppointmentService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AppointmentService) List(ctx context.Context, trainerId int64) ([]model.Appointment, error) {
	return s.repo.List(ctx, trainerId)
}

func (s *AppointmentService) Create(ctx context.Context, apt model.Appointment) (*model.Appointment, error) {
	// Run all default validation rules
	if err := apt.Validate(model.DefaultValidationRules); err != nil {
		return nil, err
	}

	// Check trainer availability
	trainerBookings, err := s.repo.GetTrainerBookings(ctx, apt.TrainerId, apt.StartTime, apt.EndTime)
	if err != nil {
		return nil, errors.InternalError("checking trainer availability: %w", err)
	}
	if len(trainerBookings) > 0 {
		errMsg := fmt.Sprintf("trainer %d is not available between %v and %v", apt.TrainerId, apt.StartTime, apt.EndTime)
		return nil, errors.ConflictError(errMsg)
	}

	// Check client availability
	clientBookings, err := s.repo.GetClientBookings(ctx, apt.UserId, apt.StartTime, apt.EndTime)
	if err != nil {
		return nil, errors.InternalError("checking user availability: %w", err)
	}
	if len(clientBookings) > 0 {
		errMsg := fmt.Sprintf("user %d is not available between %v and %v", apt.UserId, apt.StartTime, apt.EndTime)
		return nil, errors.ConflictError(errMsg)
	}

	// VALID!  Create the appointment!
	return s.repo.Create(ctx, apt)
}

func (s *AppointmentService) GetAvailability(ctx context.Context, trainerID int64, windowStartsAtUTC time.Time, windowEndsAtUTC time.Time) ([]model.TimeSlot, error) {
	// Ensure input times are UTC
	windowStartsAtUTC = windowStartsAtUTC.UTC()
	windowEndsAtUTC = windowEndsAtUTC.UTC()

	// Get all booked appointments in the time range
	booked, err := s.repo.GetTrainerBookings(ctx, trainerID, windowStartsAtUTC, windowEndsAtUTC)
	if err != nil {
		return nil, err
	}

	// Load Pacific timezone
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, fmt.Errorf("failed to load Pacific timezone: %w", err)
	}

	// Round start time up to next :00 or :30
	currentSlotStart := roundUpToNextSlot(windowStartsAtUTC)
	s.logger.Info("Slot calculation",
		"original_start", windowStartsAtUTC.Format(time.RFC3339),
		"rounded_start", currentSlotStart.Format(time.RFC3339))

	var available []model.TimeSlot
	for currentSlotStart.Add(appointmentDuration).Before(windowEndsAtUTC) ||
		currentSlotStart.Add(appointmentDuration).Equal(windowEndsAtUTC) {

		currentSlotEnd := currentSlotStart.Add(appointmentDuration)

		// Convert UTC slot time to Pacific to check business hours
		slotStartPacific := currentSlotStart.In(loc)

		// Extract hour in Pacific time
		hour := slotStartPacific.Hour()

		// Check if slot starts during business hours (8 AM - 5 PM Pacific)
		if hour >= 8 && hour < 17 {
			// Check if this slot overlaps with any booked appointments
			isAvailable := true
			for _, bookedApt := range booked {
				// A slot overlaps if it starts before the booked appointment ends
				// AND ends after the booked appointment starts
				if currentSlotStart.Before(bookedApt.EndTime) &&
					currentSlotEnd.After(bookedApt.StartTime) {
					isAvailable = false
					break
				}
			}

			if isAvailable {
				available = append(available, model.TimeSlot{
					StartTime: currentSlotStart.UTC(),
					EndTime:   currentSlotEnd.UTC(),
				})
			}
		}

		currentSlotStart = currentSlotEnd
	}

	return available, nil
}

// roundUpToNextSlot rounds up a time to the next :00 or :30 minute mark
func roundUpToNextSlot(t time.Time) time.Time {
	t = t.UTC()
	minute := t.Minute()

	// If we're already on a slot boundary (:00 or :30), don't round up
	if minute == 0 || minute == 30 {
		return t
	}

	if minute < 30 {
		// Round up to next :30
		return t.Truncate(time.Hour).Add(30 * time.Minute)
	}
	// Round up to next :00
	return t.Truncate(time.Hour).Add(time.Hour)
}
