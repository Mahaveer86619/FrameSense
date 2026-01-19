package services

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/Mahaveer86619/FrameSense/pkg/consts"
	"github.com/Mahaveer86619/FrameSense/pkg/db"
	"github.com/Mahaveer86619/FrameSense/pkg/errz"
	"github.com/Mahaveer86619/FrameSense/pkg/models"
	"github.com/Mahaveer86619/FrameSense/pkg/services/storage"
	"github.com/Mahaveer86619/FrameSense/pkg/views"
	"gorm.io/gorm"
)

type StreamingService struct {
	DB      *gorm.DB
	Storage storage.StorageService
}

func NewStreamingService(storage storage.StorageService) *StreamingService {
	return &StreamingService{
		DB:      db.GetDB(),
		Storage: storage,
	}
}

// GetPlayableVideos returns all videos that are ready for streaming
func (s *StreamingService) GetPlayableVideos() ([]views.VideoStreamResponse, error) {
	var videos []models.Video

	err := s.DB.
		Where("status = ?", consts.VideoStatusReady).
		Preload("OwnerUser").
		Order("created_at DESC").
		Find(&videos).Error

	if err != nil {
		return nil, err
	}

	response := make([]views.VideoStreamResponse, len(videos))
	for i, video := range videos {
		response[i] = views.ToVideoStreamResponse(&video)
	}

	return response, nil
}

// GetVideoDetails returns details about a specific video for streaming
func (s *StreamingService) GetVideoDetails(videoID uint) (*views.VideoStreamResponse, error) {
	var video models.Video

	err := s.DB.
		Preload("OwnerUser").
		First(&video, videoID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errz.New(errz.NotFound, "video not found", err)
		}
		return nil, errz.New(errz.InternalServerError, "failed to fetch video", err)
	}

	// Check if video is ready for streaming
	if video.Status != consts.VideoStatusReady {
		return nil, errz.New(
			errz.BadRequest,
			fmt.Sprintf("video is not ready for streaming (status: %s)", video.Status),
			nil,
		)
	}

	response := views.ToVideoStreamResponse(&video)
	return &response, nil
}

// GetMasterPlaylist returns the master HLS playlist for a video
func (s *StreamingService) GetMasterPlaylist(videoID uint) (io.ReadCloser, string, error) {
	var video models.Video

	err := s.DB.First(&video, videoID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", errz.New(errz.NotFound, "video not found", err)
		}
		return nil, "", errz.New(errz.InternalServerError, "failed to fetch video", err)
	}

	// Check if video is ready
	if video.Status != consts.VideoStatusReady {
		return nil, "", errz.New(
			errz.BadRequest,
			fmt.Sprintf("video is not ready for streaming (status: %s)", video.Status),
			nil,
		)
	}

	// Get the master playlist path
	// Expected format: s3://bucket/hls/video_1/master.m3u8
	playlistPath := fmt.Sprintf("s3://%s/hls/video_%d/master.m3u8", s.getS3Bucket(), videoID)

	// If ProcessedFilePath is set, use that instead
	if video.ProcessedFilePath != "" {
		playlistPath = video.ProcessedFilePath
	}

	reader, err := s.Storage.GetHLSPlaylist(playlistPath)
	if err != nil {
		return nil, "", errz.New(errz.NotFound, "playlist not found", err)
	}

	return reader, "application/vnd.apple.mpegurl", nil
}

// GetHLSFile returns an HLS file (segment or variant playlist)
func (s *StreamingService) GetHLSFile(videoID uint, filePath string) (io.ReadCloser, string, error) {
	var video models.Video

	err := s.DB.First(&video, videoID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", errz.New(errz.NotFound, "video not found", err)
		}
		return nil, "", errz.New(errz.InternalServerError, "failed to fetch video", err)
	}

	// Check if video is ready
	if video.Status != consts.VideoStatusReady {
		return nil, "", errz.New(
			errz.BadRequest,
			fmt.Sprintf("video is not ready for streaming (status: %s)", video.Status),
			nil,
		)
	}

	// Sanitize file path to prevent directory traversal
	cleanPath := filepath.Clean(filePath)
	if strings.Contains(cleanPath, "..") {
		return nil, "", errz.New(errz.BadRequest, "invalid file path", nil)
	}

	// Construct the full S3 path
	// e.g., s3://bucket/hls/video_1/segment0.ts
	fullPath := fmt.Sprintf("s3://%s/hls/video_%d/%s", s.getS3Bucket(), videoID, cleanPath)

	reader, err := s.Storage.GetHLSPlaylist(fullPath)
	if err != nil {
		return nil, "", errz.New(errz.NotFound, "file not found", err)
	}

	// Determine content type based on file extension
	contentType := getContentTypeFromPath(filePath)

	return reader, contentType, nil
}

// Helper function to get S3 bucket name from config
func (s *StreamingService) getS3Bucket() string {
	// You'll need to add this to your config
	// For now, return a default
	return "framesense-videos"
}

// Helper function to determine content type
func getContentTypeFromPath(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	case ".ts":
		return "video/MP2T"
	case ".mp4":
		return "video/mp4"
	case ".vtt":
		return "text/vtt"
	default:
		return "application/octet-stream"
	}
}
