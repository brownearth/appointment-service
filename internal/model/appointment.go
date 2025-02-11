package model

import (
	"appointment-service/internal/errors"
	"fmt"
	"time"
)

// Appointment represents a scheduled meeting between a user and a trainer.
type Appointment struct {
	Id        int64
	StartTime time.Time
	EndTime   time.Time
	TrainerId int64
	UserId    int64
}

// Defines a type for validation rules, then we can pass
// sets of rules (as functions) to a validator method
type ValidationRule func(a *Appointment) error

// MustBeThirtyMinutes checks if the duration of the given appointment is exactly 30 minutes.
func MustBeThirtyMinutes(a *Appointment) error {
	duration := a.EndTime.Sub(a.StartTime)
	if duration != 30*time.Minute {
		return errors.ValidationError(
			fmt.Sprintf("appointment must be exactly 30 minutes, got %v", duration),
		)
	}
	return nil
}

// MustBeDuringBusinessHours checks if the appointment's start and end times
// fall within business hours (8am to 5pm Pacific Time).
func MustBeDuringBusinessHours(a *Appointment) error {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	startLocal := a.StartTime.In(loc)
	endLocal := a.EndTime.In(loc)

	year, month, day := startLocal.Date()
	businessStart := time.Date(year, month, day, 8, 0, 0, 0, loc)
	businessEnd := time.Date(year, month, day, 17, 0, 0, 0, loc)

	if startLocal.Before(businessStart) || startLocal.After(businessEnd) {
		return errors.ValidationError("appointment must start between 8am and 5pm Pacific")
	}

	if endLocal.Before(businessStart) || endLocal.After(businessEnd) {
		return errors.ValidationError("appointment must end between 8am and 5pm Pacific")
	}

	return nil
}

// DefaultValidationRules is a set of validation rules that can be used
// should you not want to define your own set of rules.
var DefaultValidationRules = []ValidationRule{
	MustBeThirtyMinutes,
	MustBeDuringBusinessHours,
}

// Validate runs the given validation rules on the appointment.
// Allows callers/clients to define the set of rules appropriate for the context.
func (a *Appointment) Validate(rules []ValidationRule) error {
	for _, validationRule := range rules {
		if err := validationRule(a); err != nil {
			return err
		}
	}
	return nil
}
