package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Mahaveer86619/FrameSense/pkg/errz"
	"github.com/Mahaveer86619/FrameSense/pkg/services"
	"github.com/Mahaveer86619/FrameSense/pkg/views"
	"github.com/labstack/echo/v4"
)

type StreamingHandler struct {
	streamingService *services.StreamingService
}

func NewStreamingHandler(g *echo.Group, service *services.StreamingService) *StreamingHandler {
	handler := &StreamingHandler{
		streamingService: service,
	}

	g.GET("/videos", handler.GetAllPlayableVideos)
	g.GET("/videos/:id", handler.GetVideoDetails)
	g.GET("/videos/:id/master.m3u8", handler.StreamMasterPlaylist)
	g.GET("/videos/:id/hls/*", handler.StreamHLSSegment)

	return handler
}

func (h *StreamingHandler) GetAllPlayableVideos(c echo.Context) error {
	videos, err := h.streamingService.GetPlayableVideos()
	if err != nil {
		errz.HandleErrors(c, errz.New(errz.InternalServerError, "failed to fetch videos", err))
		return nil
	}

	resp := &views.Success{
		StatusCode: http.StatusOK,
		Message:    "Playable videos fetched successfully",
		Data:       videos,
	}
	return resp.Send(c)
}

func (h *StreamingHandler) GetVideoDetails(c echo.Context) error {
	videoID, err := getVideoIDFromPath(c)
	if err != nil {
		errz.HandleErrors(c, errz.New(errz.BadRequest, "invalid video ID", err))
		return nil
	}

	video, err := h.streamingService.GetVideoDetails(videoID)
	if err != nil {
		errz.HandleErrors(c, err)
		return nil
	}

	resp := &views.Success{
		StatusCode: http.StatusOK,
		Message:    "Video details fetched successfully",
		Data:       video,
	}
	return resp.Send(c)
}

func (h *StreamingHandler) StreamMasterPlaylist(c echo.Context) error {
	videoID, err := getVideoIDFromPath(c)
	if err != nil {
		return errz.New(errz.BadRequest, "invalid video ID", err)
	}

	playlistReader, contentType, err := h.streamingService.GetMasterPlaylist(videoID)
	if err != nil {
		return errz.New(errz.NotFound, "playlist not found", err)
	}
	defer playlistReader.Close()

	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Cache-Control", "no-cache")

	// Stream the playlist
	_, err = io.Copy(c.Response().Writer, playlistReader)
	if err != nil {
		return errz.New(errz.InternalServerError, "failed to stream playlist", err)
	}

	return nil
}

func (h *StreamingHandler) StreamHLSSegment(c echo.Context) error {
	videoID, err := getVideoIDFromPath(c)
	if err != nil {
		return errz.New(errz.BadRequest, "invalid video ID", err)
	}

	// Get the file path from the wildcard parameter
	// e.g., /stream/videos/1/hls/segment0.ts -> segment0.ts
	// e.g., /stream/videos/1/hls/720p/playlist.m3u8 -> 720p/playlist.m3u8
	filePath := c.Param("*")
	if filePath == "" {
		return errz.New(errz.BadRequest, "file path is required", nil)
	}

	fileReader, contentType, err := h.streamingService.GetHLSFile(videoID, filePath)
	if err != nil {
		return errz.New(errz.NotFound, "file not found", err)
	}
	defer fileReader.Close()

	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	if strings.HasSuffix(filePath, ".ts") {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
	} else {
		c.Response().Header().Set("Cache-Control", "no-cache")
	}

	_, err = io.Copy(c.Response().Writer, fileReader)
	if err != nil {
		return errz.New(errz.InternalServerError, "failed to stream file", err)
	}

	return nil
}

func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	case ".ts":
		return "video/MP2T"
	case ".mp4":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}
