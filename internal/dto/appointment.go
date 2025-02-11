package dto

import (
	"encoding/json"
	"time"
)

// Request DTO Types
type CreateAppointmentRequest struct {
	TrainerId int64     `json:"trainer_id" binding:"required,gt=0"`
	StartTime time.Time `json:"start_time" binding:"required" time_format:"2006-01-02T15:04:05Z"`
	EndTime   time.Time `json:"end_time" binding:"required,gtfield=StartTime" time_format:"2006-01-02T15:04:05Z"`
	UserId    int64     `json:"user_id" binding:"required,gt=0"`
}

type ListAppointmentsRequest struct {
	TrainerId int64 `uri:"trainer_id" binding:"required"`
}

type GetAvailabilityRequest struct {
	TrainerId int64     `uri:"trainer_id"`
	StartsAt  time.Time `form:"starts_at" time_format:"2006-01-02T15:04:05Z07:00"`
	EndsAt    time.Time `form:"ends_at" time_format:"2006-01-02T15:04:05Z07:00"`
}

// Response DTO Types
type AppointmentResponse struct {
	Id        int64     `json:"id"`
	TrainerId int64     `json:"trainer_id"`
	StartTime time.Time `json:"start_time" time_format:"2006-01-02T15:04:05Z"`
	EndTime   time.Time `json:"end_time" time_format:"2006-01-02T15:04:05Z"`
	UserId    int64     `json:"user_id"`
}

type AvailabilityResponse struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// Custom marshaler for AvailabilityResponse to ensure UTC output
func (r AvailabilityResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}{
		StartTime: r.StartTime.UTC().Format(time.RFC3339),
		EndTime:   r.EndTime.UTC().Format(time.RFC3339),
	})
}

// Ensure times are in UTC after binding
func (r *GetAvailabilityRequest) AfterBinding() error {
	r.StartsAt = r.StartsAt.UTC()
	r.EndsAt = r.EndsAt.UTC()
	return nil
}
