package services

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/Mahaveer86619/FrameSense/pkg/consts"
	"github.com/Mahaveer86619/FrameSense/pkg/db"
	"github.com/Mahaveer86619/FrameSense/pkg/models"
	"github.com/Mahaveer86619/FrameSense/pkg/services/queue"
	"github.com/Mahaveer86619/FrameSense/pkg/services/storage"
	"gorm.io/gorm"
)

type VideoProcessingService struct {
	DB      *gorm.DB
	Storage storage.StorageService
	Queue   queue.QueueService
}

func NewVideoProcessingService(
	storage storage.StorageService,
	queueService queue.QueueService,
) *VideoProcessingService {
	return &VideoProcessingService{
		DB:      db.GetDB(),
		Storage: storage,
		Queue:   queueService,
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
		_ = s.Storage.DeleteVideo(sourcePath)
		return nil, fmt.Errorf("failed to create video record: %w", err)
	}

	urlExpiration := 24 * time.Hour

	downloadURL, err := s.Storage.GeneratePresignedDownloadURL(sourcePath, urlExpiration)
	if err != nil {
		s.updateVideoStatus(video.ID, string(consts.VideoStatusFailed))
		return nil, fmt.Errorf("failed to generate download URL: %w", err)
	}

	processedFilename := fmt.Sprintf("processed_%d_%d_%s", ownerID, time.Now().Unix(), originalFilename)
	uploadURL, err := s.Storage.GeneratePresignedUploadURL(processedFilename, urlExpiration)
	if err != nil {
		s.updateVideoStatus(video.ID, string(consts.VideoStatusFailed))
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	processedPath := fmt.Sprintf("s3://%s/processed/%s", "your-bucket", processedFilename)

	ingestMsg := models.NewVideoIngestMessage(
		video.ID,
		ownerID,
		title,
		downloadURL,
		uploadURL,
		processedPath,
		time.Now().Add(urlExpiration),
	)

	if err := s.Queue.SendVideoIngestMessage(ingestMsg); err != nil {
		s.updateVideoStatus(video.ID, string(consts.VideoStatusQueueFailed))
		return nil, fmt.Errorf("video uploaded but failed to queue for processing: %w", err)
	}

	video.Status = consts.VideoStatusQueued
	s.DB.Save(video)

	return video, nil
}

func (s *VideoProcessingService) updateVideoStatus(videoID uint, status string) {
	s.DB.Model(&models.Video{}).Where("id = ?", videoID).Update("status", status)
}

func (s *VideoProcessingService) HandleProcessingComplete(
	videoID uint,
	processedPath string,
	success bool,
	errorMsg string,
) error {
	video := &models.Video{}
	if err := s.DB.First(video, videoID).Error; err != nil {
		return fmt.Errorf("video not found: %w", err)
	}

	if success {
		video.ProcessedFilePath = processedPath
		video.Status = consts.VideoStatusReady
	} else {
		video.Status = consts.VideoStatusFailed
		video.ErrorMessage = errorMsg
	}

	return s.DB.Save(video).Error
}
