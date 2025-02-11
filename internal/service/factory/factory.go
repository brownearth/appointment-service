package factory

import (
	"appointment-service/internal/repository"
	"appointment-service/internal/service"
	"log/slog"
)

// NewAppointmentService creates a new appointment service with all its dependencies
// Dont really need a factory for this, as there is only one
// but it's here for consistency
func NewAppointmentService(repo repository.AppointmentRepository, logger *slog.Logger) service.AppointmentServicer {
	return service.NewAppointmentService(repo, logger.With("service", "AppointmentService"))
}
