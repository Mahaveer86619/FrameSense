package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/FrameSense/pkg/errz"
	"github.com/Mahaveer86619/FrameSense/pkg/services"
	"github.com/Mahaveer86619/FrameSense/pkg/views"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(g *echo.Group, service *services.AuthService) *AuthHandler {
	handler := &AuthHandler{service: service}

	g.POST("/register", handler.Register)
	g.POST("/login", handler.Login)
	g.POST("/refresh", handler.Refresh)

	return handler
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req views.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return errz.New(errz.BadRequest, "Invalid request body", err)
	}

	// Add manual validation here if not using a validator middleware
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return errz.New(errz.BadRequest, "Missing required fields", nil)
	}

	resp, err := h.service.Register(req)
	if err != nil {
		errz.HandleErrors(c, err)
		return nil
	}

	return c.JSON(http.StatusCreated, views.Success{
		StatusCode: http.StatusCreated,
		Message:    "Registration successful",
		Data:       resp,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req views.LoginRequest
	if err := c.Bind(&req); err != nil {
		return errz.New(errz.BadRequest, "Invalid request body", err)
	}

	resp, err := h.service.Login(req)
	if err != nil {
		errz.HandleErrors(c, err)
		return nil
	}

	return c.JSON(http.StatusOK, views.Success{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data:       resp,
	})
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var req views.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return errz.New(errz.BadRequest, "Invalid request body", err)
	}

	tokens, err := h.service.Refresh(req)
	if err != nil {
		errz.HandleErrors(c, err)
		return nil
	}

	return c.JSON(http.StatusOK, views.Success{
		StatusCode: http.StatusOK,
		Message:    "Tokens refreshed",
		Data:       tokens,
	})
}
