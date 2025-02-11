package api

import (
	"appointment-service/internal/dto"
	"appointment-service/internal/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ListAppointments is a handler to list all appointments for a given trainer
// Attach this function to Server struct for ease of setting up routes
func (s *Server) ListAppointments(c *gin.Context) {

	// Bind parameters to ListAppointmentsRequest DTO
	// -----------------------------------------------
	var req dto.ListAppointmentsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		handleError(c, errors.ValidationError(err.Error()))
		return
	}

	// List the appointments
	// ----------------------
	appointments, err := s.appointmentService.List(c.Request.Context(), req.TrainerId)
	if err != nil {
		handleError(c, err)
		return
	}

	// Convert appointments to AppointmentResponse DTOs
	// ------------------------------------------------
	response := dto.ToListAppointmentsResponse(appointments)
	c.JSON(http.StatusOK, response)
}

// CreateAppointment is a handler to create a new appointment
// Attach this function to Server struct for ease of setting up routes
func (s *Server) CreateAppointment(c *gin.Context) {

	// Bind parameters to CreateAppointmentRequest DTO
	// -----------------------------------------------
	var req dto.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.ValidationError(err.Error()))
		return
	}

	// Convert the request to a model
	//
	appointment := dto.ToAppointmentModel(&req)

	// Create the appointment
	// -----------------------
	createdAppointment, err := s.appointmentService.Create(c.Request.Context(), appointment)
	if err != nil {
		handleError(c, err)
		return
	}

	// Convert the created appointment to a response DTO
	// -------------------------------------------------
	response := dto.ToAppointmentResponse(createdAppointment)
	c.JSON(http.StatusCreated, response)
}

// GetAvailability is a handler to get available slots for a given trainer
func (s *Server) GetAvailability(c *gin.Context) {

	// This was a little wonky, as there are URL parameters and query parameters
	// to bind.  We'll bind the URL parameters first, then the query parameters.
	// But we could not use the required annotation on the DTO struct
	// so we have to validate explicitly.
	var req dto.GetAvailabilityRequest
	if err := c.ShouldBindUri(&req); err != nil {
		s.logger.Error("URI binding failed", "error", err)
		handleError(c, errors.ValidationError(err.Error()))
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		s.logger.Error("Query binding failed",
			"error", err,
			"raw_starts_at", c.Query("starts_at"),
			"raw_ends_at", c.Query("ends_at"))
		handleError(c, errors.ValidationError(err.Error()))
		return
	}

	// Ensure times are in UTC
	// -----------------------
	req.StartsAt = req.StartsAt.UTC()
	req.EndsAt = req.EndsAt.UTC()

	s.logger.Info("Request times in UTC",
		"starts_at", req.StartsAt.Format(time.RFC3339),
		"ends_at", req.EndsAt.Format(time.RFC3339))

	// Validate the request
	// --------------------
	if err := validateAvailabilityRequest(&req); err != nil {
		handleError(c, err)
		return
	}

	// Get available slots
	// -------------------
	available, err := s.appointmentService.GetAvailability(
		c.Request.Context(),
		req.TrainerId,
		req.StartsAt,
		req.EndsAt,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	// Convert to response DTOs, ensuring UTC
	// --------------------------------------
	response := make([]dto.AvailabilityResponse, len(available))
	for i, slot := range available {
		response[i] = dto.AvailabilityResponse{
			StartTime: slot.StartTime.UTC(),
			EndTime:   slot.EndTime.UTC(),
		}
	}

	c.JSON(http.StatusOK, response)
}

func validateAvailabilityRequest(req *dto.GetAvailabilityRequest) error {

	if req.TrainerId <= 0 {
		return errors.ValidationError("trainer_id must be greater than 0")
	}

	if req.StartsAt.IsZero() {
		return errors.ValidationError("starts_at is required and must be a valid timestamp")
	}

	if req.EndsAt.IsZero() {
		return errors.ValidationError("ends_at is required and must be a valid timestamp")
	}

	if req.EndsAt.Before(req.StartsAt) {
		return errors.ValidationError("ends_at must be after starts_at")
	}

	return nil
}

// handleError is a helper function to handle different types of errors
func handleError(c *gin.Context, err error) {
	// If this is an application error, return the error message and status code from within
	if appErr, ok := errors.IsAppError(err); ok {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}

	// For any other error, return 500
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
