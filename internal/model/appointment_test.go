package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMustBeDuringBusinessHours tests the validation of appointment times against business hours.
//
// It includes the following test cases:
//
// * Valid appointment within business hours (9:00 AM - 5:00 PM)
// * Appointment starting before business hours
// * Appointment ending after business hours
// * Appointment completely outside business hours
//
// Note: Business hours are assumed to be 9:00 AM to 5:00 PM Pacific Time.
// All test cases use the America/Los_Angeles timezone to ensure consistent validation.
func TestMustBeDuringBusinessHours(t *testing.T) {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		wantErr   bool
	}{
		{
			name:      "valid appointment within business hours",
			startTime: time.Date(2023, 10, 10, 9, 0, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 9, 30, 0, 0, loc),
			wantErr:   false,
		},
		{
			name:      "appointment starts before business hours",
			startTime: time.Date(2023, 10, 10, 7, 30, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 8, 0, 0, 0, loc),
			wantErr:   true,
		},
		{
			name:      "appointment ends after business hours",
			startTime: time.Date(2023, 10, 10, 16, 30, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 17, 30, 0, 0, loc),
			wantErr:   true,
		},
		{
			name:      "appointment starts and ends outside business hours",
			startTime: time.Date(2023, 10, 10, 18, 0, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 18, 30, 0, 0, loc),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appointment := &Appointment{
				StartTime: tt.startTime,
				EndTime:   tt.endTime,
			}
			err := MustBeDuringBusinessHours(appointment)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestMustBeThirtyMinutes tests the validation of appointment duration.
//
// It includes the following test cases:
//
// * Valid appointment of exactly 30 minutes
// * Appointment less than 30 minutes
// * Appointment more than 30 minutes
//
// Note: This validation ensures that all appointments are exactly 30 minutes
// in duration, which is a business requirement. The test uses fixed times
// in the America/Los_Angeles timezone for consistency.
func TestMustBeThirtyMinutes(t *testing.T) {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		wantErr   bool
	}{
		{
			name:      "valid 30 minute appointment",
			startTime: time.Date(2023, 10, 10, 9, 0, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 9, 30, 0, 0, loc),
			wantErr:   false,
		},
		{
			name:      "appointment less than 30 minutes",
			startTime: time.Date(2023, 10, 10, 9, 0, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 9, 20, 0, 0, loc),
			wantErr:   true,
		},
		{
			name:      "appointment more than 30 minutes",
			startTime: time.Date(2023, 10, 10, 9, 0, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 9, 40, 0, 0, loc),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appointment := &Appointment{
				StartTime: tt.startTime,
				EndTime:   tt.endTime,
			}
			err := MustBeThirtyMinutes(appointment)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidate tests the complete validation pipeline for appointments.
//
// It includes the following test cases:
//
// * Valid appointment meeting all validation rules
// * Appointment outside business hours
// * Appointment with incorrect duration
// * Appointment with multiple validation failures
// * Appointment with no validation rules applied
//
// Note: The DefaultValidationRules are used for most test cases to ensure
// that appointments meet all standard business requirements. One test case
// uses an empty rule set to verify behavior when no validation is required.
// All times are in America/Los_Angeles timezone for consistency.
func TestValidate(t *testing.T) {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		rules     []ValidationRule
		wantErr   bool
	}{
		{
			name:      "valid appointment with all rules",
			startTime: time.Date(2023, 10, 10, 9, 0, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 9, 30, 0, 0, loc),
			rules:     DefaultValidationRules,
			wantErr:   false,
		},
		{
			name:      "appointment starts before business hours",
			startTime: time.Date(2023, 10, 10, 7, 30, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 8, 0, 0, 0, loc),
			rules:     DefaultValidationRules,
			wantErr:   true,
		},
		{
			name:      "appointment less than 30 minutes",
			startTime: time.Date(2023, 10, 10, 9, 0, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 9, 20, 0, 0, loc),
			rules:     DefaultValidationRules,
			wantErr:   true,
		},
		{
			name:      "appointment more than 30 minutes",
			startTime: time.Date(2023, 10, 10, 9, 0, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 9, 40, 0, 0, loc),
			rules:     DefaultValidationRules,
			wantErr:   true,
		},
		{
			name:      "valid appointment with no rules",
			startTime: time.Date(2023, 10, 10, 9, 0, 0, 0, loc),
			endTime:   time.Date(2023, 10, 10, 9, 30, 0, 0, loc),
			rules:     []ValidationRule{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appointment := &Appointment{
				StartTime: tt.startTime,
				EndTime:   tt.endTime,
			}
			err := appointment.Validate(tt.rules)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
