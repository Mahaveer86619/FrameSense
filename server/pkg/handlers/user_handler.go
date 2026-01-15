package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/FrameSense/pkg/middleware"
	"github.com/Mahaveer86619/FrameSense/pkg/services"
	"github.com/Mahaveer86619/FrameSense/pkg/utils"
	"github.com/Mahaveer86619/FrameSense/pkg/views"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(rg *echo.Group, userService *services.UserService) *UserHandler {
	handler := &UserHandler{
		userService: userService,
	}

	users := rg.Group("/users")

	users.GET("/me", handler.GetUserByToken)
	users.GET("/:id", handler.GetUserByID)

	return handler
}

func (h *UserHandler) GetUserByToken(c echo.Context) error {
	currentUserId := middleware.GetUserID(c)
	if currentUserId == 0 {
		return echo.ErrUnauthorized
	}

	user, err := h.userService.FindByID(currentUserId)
	if err != nil {
		return err
	}

	resp := &views.Success{
		StatusCode: http.StatusOK,
		Message:    "User fetched successfully",
		Data:       views.ToUserResponse(user),
	}
	return resp.Send(c)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	maskedId := c.Param("id")
	if maskedId == "" {
		return echo.ErrBadRequest
	}

	id, err := utils.UnmaskID(maskedId)
	if err != nil {
		return echo.ErrBadRequest
	}

	user, err := h.userService.FindByID(id)
	if err != nil {
		return err
	}

	resp := &views.Success{
		StatusCode: http.StatusOK,
		Message:    "User fetched successfully",
		Data:       views.ToUserResponse(user),
	}
	return resp.Send(c)
}
