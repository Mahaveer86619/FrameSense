package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/FrameSense/pkg/errz"
	"github.com/Mahaveer86619/FrameSense/pkg/middleware"
	"github.com/Mahaveer86619/FrameSense/pkg/services"
	"github.com/Mahaveer86619/FrameSense/pkg/views"
	"github.com/labstack/echo/v4"
)

type VideoProcessingHandler struct {
	Service *services.VideoProcessingService
}

func NewVideoProcessingHandler(
	g *echo.Group,
	service *services.VideoProcessingService,
) *VideoProcessingHandler {
	h := &VideoProcessingHandler{
		Service: service,
	}
	g.POST("/videos", h.UploadVideo)
	return h
}

func (h *VideoProcessingHandler) UploadVideo(c echo.Context) error {
	var req views.CreateVideoRequest
	title := c.FormValue("title")
	if title == "" {
		return errz.New(errz.BadRequest, "title is required", nil)
	}
	req.Title = title

	description := c.FormValue("description")
	req.Description = description

	fileHeader, err := c.FormFile("video")
	if err != nil {
		return errz.New(errz.BadRequest, "video file is required", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return errz.New(errz.InternalServerError, "failed to open file", err)
	}
	defer file.Close()

	// From auth middleware
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	video, err := h.Service.UploadVideo(
		userID,
		req.Title,
		req.Description,
		file,
		fileHeader.Filename,
	)
	if err != nil {
		errz.HandleErrors(c, errz.New(errz.InternalServerError, "failed to upload video", err))
		return nil
	}

	resp := &views.Success{
		StatusCode: http.StatusCreated,
		Message:    "Video uploaded successfully",
		Data:       views.ToVideoResponse(video),
	}
	return resp.Send(c)
}
