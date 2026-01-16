package models

import (
	"time"

	"github.com/google/uuid"
)

type VideoIngestMessage struct {
	JobID       string              `json:"job_id"`
	VideoID     uint                `json:"video_id"`
	DownloadURL string              `json:"download_url"`
	Callback    VideoIngestCallback `json:"callback"`
}

type VideoIngestCallback struct {
	StatusURL        string `json:"status_url"`         // POST status updates
	UploadURL        string `json:"upload_url"`         // POST HLS files directly
	RequestUploadURL string `json:"request_upload_url"` // GET presigned URLs
}

func NewVideoIngestMessage(
	videoID uint,
	downloadURL string,
	callbackStatusURL string,
	callbackUploadURL string,
	callbackRequestUploadURL string,
) *VideoIngestMessage {
	return &VideoIngestMessage{
		JobID:       uuid.New().String(),
		VideoID:     videoID,
		DownloadURL: downloadURL,
		Callback: VideoIngestCallback{
			StatusURL:        callbackStatusURL,
			UploadURL:        callbackUploadURL,
			RequestUploadURL: callbackRequestUploadURL,
		},
	}
}

type VideoProcessingRetry struct {
	VideoID      uint      `json:"video_id"`
	AttemptCount int       `json:"attempt_count"`
	LastError    string    `json:"last_error"`
	RetryAt      time.Time `json:"retry_at"`
}
