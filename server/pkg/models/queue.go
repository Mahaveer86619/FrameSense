package models

import "time"

type VideoIngestMessage struct {
	VideoID            uint      `json:"video_id"`
	OwnerID            uint      `json:"owner_id"`
	Title              string    `json:"title"`
	SourceDownloadURL  string    `json:"source_download_url"`  // Presigned download URL
	ProcessedUploadURL string    `json:"processed_upload_url"` // Presigned upload URL
	ProcessedFilePath  string    `json:"processed_file_path"`  // Expected S3 path after upload
	ExpiresAt          time.Time `json:"expires_at"`
}

func NewVideoIngestMessage(
	videoID uint,
	ownerID uint,
	title string,
	sourceDownloadURL string,
	processedUploadURL string,
	processedFilePath string,
	expiresAt time.Time,
) *VideoIngestMessage {
	return &VideoIngestMessage{
		VideoID:            videoID,
		OwnerID:            ownerID,
		Title:              title,
		SourceDownloadURL:  sourceDownloadURL,
		ProcessedUploadURL: processedUploadURL,
		ProcessedFilePath:  processedFilePath,
		ExpiresAt:          expiresAt,
	}
}

type VideoProcessingRetry struct {
	VideoID      uint      `json:"video_id"`
	AttemptCount int       `json:"attempt_count"`
	LastError    string    `json:"last_error"`
	RetryAt      time.Time `json:"retry_at"`
}
