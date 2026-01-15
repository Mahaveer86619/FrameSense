package services

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/Mahaveer86619/FrameSense/pkg/consts"
	"github.com/Mahaveer86619/FrameSense/pkg/db"
	"github.com/Mahaveer86619/FrameSense/pkg/models"
	"github.com/Mahaveer86619/FrameSense/pkg/services/storage"
	"gorm.io/gorm"
)

type VideoProcessingService struct {
	DB      *gorm.DB
	Storage storage.StorageService
}

func NewVideoProcessingService(
	storage storage.StorageService,
) *VideoProcessingService {
	return &VideoProcessingService{
		DB:      db.GetDB(),
		Storage: storage,
	}
}

func (s *VideoProcessingService) UploadVideo(
	ownerID uint,
	title string,
	description string,
	file multipart.File,
	originalFilename string,
) (*models.Video, error) {

	filename := fmt.Sprintf(
		"video_%d_%d_%s",
		ownerID,
		time.Now().Unix(),
		originalFilename,
	)

	sourcePath, err := s.Storage.SaveVideo(file, filename)
	if err != nil {
		return nil, err
	}

	video := &models.Video{
		Title:          title,
		Description:    description,
		SourceFilePath: sourcePath,
		OwnerUserId:    ownerID,
		Status:         consts.VideoStatusUploaded,
	}

	if err := s.DB.Create(video).Error; err != nil {
		return nil, err
	}

	// NEXT STEP (later):
	// Push SQS message: video.ingest

	return video, nil
}
