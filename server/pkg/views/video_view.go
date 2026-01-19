package views

import (
	"fmt"
	"time"

	"github.com/Mahaveer86619/FrameSense/pkg/consts"
	"github.com/Mahaveer86619/FrameSense/pkg/models"
	"github.com/Mahaveer86619/FrameSense/pkg/utils"
)

type CreateVideoRequest struct {
	Title       string `form:"title" validate:"required,min=3,max=200"`
	Description string `form:"description"`
}

type VideoResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      consts.VideoStatus `json:"status"`
	Duration    string             `json:"duration"`
	CreatedAt   time.Time          `json:"created_at"`
}

func ToVideoResponse(v *models.Video) VideoResponse {
	return VideoResponse{
		ID:          utils.MaskID(v.ID),
		Title:       v.Title,
		Description: v.Description,
		Status:      v.Status,
		Duration:    v.Duration,
		CreatedAt:   v.CreatedAt,
	}
}

type VideoStreamResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Duration     string    `json:"duration,omitempty"`
	Status       string    `json:"status"`
	Owner        UserInfo  `json:"owner"`
	StreamURL    string    `json:"stream_url"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type VideoListResponse struct {
	Videos     []VideoStreamResponse `json:"videos"`
	TotalCount int                   `json:"total_count"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
}

func ToVideoStreamResponse(video *models.Video) VideoStreamResponse {
	baseURL := getBaseURL()

	return VideoStreamResponse{
		ID:          video.ID,
		Title:       video.Title,
		Description: video.Description,
		Duration:    video.Duration,
		Status:      string(video.Status),
		Owner: UserInfo{
			ID:       video.OwnerUser.ID,
			Username: video.OwnerUser.Username,
		},
		StreamURL:    fmt.Sprintf("%s/api/stream/videos/%d/master.m3u8", baseURL, video.ID),
		ThumbnailURL: fmt.Sprintf("%s/api/stream/videos/%d/thumbnail.jpg", baseURL, video.ID),
		CreatedAt:    video.CreatedAt,
		UpdatedAt:    video.UpdatedAt,
	}
}

func getBaseURL() string {
	return "http://localhost:7000"
}
