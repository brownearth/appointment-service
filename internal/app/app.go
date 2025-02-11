package app

import (
	"appointment-service/internal/api"
	"appointment-service/internal/config"
	"appointment-service/internal/repository"
	repofactory "appointment-service/internal/repository/factory"
	"appointment-service/internal/service"
	servicefactory "appointment-service/internal/service/factory"
	"log/slog"
)

// Application contains all dependencies
type Application struct {
	Config             *config.Config
	Logger             *slog.Logger
	Repository         repository.AppointmentRepository
	AppointmentService service.AppointmentServicer
	Server             *api.Server
}

// New creates a new application instance with all dependencies wired up
func New(cfg *config.Config, logger *slog.Logger) (*Application, error) {

	// Create repository
	// ------------------
	repo, err := repofactory.NewRepository(cfg, logger)
	if err != nil {
		return nil, err
	}

	// Create service, injecting the repository
	// ----------------------------------------
	appointmentService := servicefactory.NewAppointmentService(repo, logger)

	// Create server
	// -------------
	server, err := api.NewServer(cfg, appointmentService, logger)
	if err != nil {
		return nil, err
	}

	return &Application{
		Config:             cfg,
		Logger:             logger,
		Repository:         repo,
		AppointmentService: appointmentService,
		Server:             server,
	}, nil
}

// Close cleans up application resources
func (app *Application) Close() error {
	return app.Repository.Close()
}
