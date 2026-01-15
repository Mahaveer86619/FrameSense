package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/FrameSense/pkg/services"
	"github.com/Mahaveer86619/FrameSense/pkg/views"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	s *services.HealthService
}

func NewHealthHandler(g *echo.Group, service *services.HealthService) *HealthHandler {
	h := &HealthHandler{
		s: service,
	}

	g.GET("/health", h.HandleHealthCheck)

	return h
}

func (hh *HealthHandler) HandleHealthCheck(c echo.Context) error {
	msg, ok := hh.s.CheckHealth()

	if ok {
		data := views.NewHealthView(msg)

		resp := &views.Success{
			StatusCode: http.StatusOK,
			Message:    "System is operational",
			Data:       data,
		}
		return resp.Send(c)
	}

	resp := &views.Failure{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "System is unhealthy",
	}
	return resp.Send(c)
}
