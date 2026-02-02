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
	g.POST("/videos/:id/status", h.UpdateStatus)
	g.POST("/videos/:id/hls/upload", h.UploadHLSFile)
	g.GET("/videos/:id/hls/request-upload", h.RequestHLSUploadURL)
	return h
}

func (h *VideoProcessingHandler) UploadVideo(c echo.Context) error {
	var req views.CreateVideoRequest
	title := c.FormValue("title")
	if title == "" {
		errz.HandleErrors(c, errz.New(errz.BadRequest, "title is required", nil))
		return nil
	}
	req.Title = title

	description := c.FormValue("description")
	req.Description = description

	fileHeader, err := c.FormFile("video")
	if err != nil {
		errz.HandleErrors(c, errz.New(errz.BadRequest, "video file is required", err))
		return nil
	}

	file, err := fileHeader.Open()
	if err != nil {
		errz.HandleErrors(c, errz.New(errz.InternalServerError, "failed to open file", err))
		return nil
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

func (h *VideoProcessingHandler) UpdateStatus(c echo.Context) error {
	videoID, err := getVideoIDFromPath(c)
	if err != nil {
		return errz.New(errz.BadRequest, "invalid video ID", err)
	}

	var req struct {
		Status       string `json:"status"`
		ErrorMessage string `json:"error_message,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		return errz.New(errz.BadRequest, "invalid request body", err)
	}

	if err := h.Service.UpdateVideoStatus(videoID, req.Status, req.ErrorMessage); err != nil {
		return errz.New(errz.InternalServerError, "failed to update status", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "status updated successfully",
	})
}

func (h *VideoProcessingHandler) RequestHLSUploadURL(c echo.Context) error {
	videoID, err := getVideoIDFromPath(c)
	if err != nil {
		return errz.New(errz.BadRequest, "invalid video ID", err)
	}

	filename := c.QueryParam("filename")
	if filename == "" {
		return errz.New(errz.BadRequest, "filename query parameter is required", nil)
	}

	uploadURL, err := h.Service.GenerateHLSUploadURL(videoID, filename)
	if err != nil {
		return errz.New(errz.InternalServerError, "failed to generate upload URL", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"upload_url": uploadURL,
	})
}

func (h *VideoProcessingHandler) UploadHLSFile(c echo.Context) error {
	videoID, err := getVideoIDFromPath(c)
	if err != nil {
		return errz.New(errz.BadRequest, "invalid video ID", err)
	}

	filename := c.FormValue("filename")
	if filename == "" {
		return errz.New(errz.BadRequest, "filename is required", nil)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return errz.New(errz.BadRequest, "file is required", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return errz.New(errz.InternalServerError, "failed to open file", err)
	}
	defer file.Close()

	if err := h.Service.ReceiveHLSFile(videoID, filename, file); err != nil {
		return errz.New(errz.InternalServerError, "failed to save HLS file", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "file uploaded successfully",
	})
}

func getVideoIDFromPath(c echo.Context) (uint, error) {
	var videoID uint
	if err := echo.PathParamsBinder(c).Uint("id", &videoID).BindError(); err != nil {
		return 0, err
	}
	return videoID, nil
}
