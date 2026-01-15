package views

import (
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
