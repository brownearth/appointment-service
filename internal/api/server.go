package api

import (
	"appointment-service/internal/config"
	"appointment-service/internal/middleware"
	"appointment-service/internal/service"

	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer         *http.Server
	router             *gin.Engine
	cfg                *config.Config
	appointmentService service.AppointmentServicer
	logger             *slog.Logger
}

// NewServer creates a new instance of the server
func NewServer(cfg *config.Config, appointmentService service.AppointmentServicer, logger *slog.Logger) (*Server, error) {

	r := gin.New()

	server := &Server{
		httpServer:         &http.Server{},
		router:             r,
		cfg:                cfg,
		appointmentService: appointmentService,
		logger:             logger,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server, nil
}

// setupMiddleware configures the server's middleware
func (s *Server) setupMiddleware() {
	s.router.Use(middleware.GinLogger(s.logger))
	s.router.Use(gin.Recovery()) // <-- panic to 500 conversion
}

// setupRoutes configures the server's routes
func (s *Server) setupRoutes() {

	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/appointments/trainers/:trainer_id", s.ListAppointments)
		v1.POST("/appointments", s.CreateAppointment)
		v1.GET("/appointments/trainers/:trainer_id/availability", s.GetAvailability)
	}
}

func (s *Server) Run(addr string) error {
	s.httpServer.Addr = addr
	s.httpServer.Handler = s.router
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
