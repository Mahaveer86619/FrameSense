package queue

// VideoIngestMessage represents the strict contract from the Queue.
// It contains all necessary metadata so the worker remains stateless.
type VideoIngestMessage struct {
	JobID       string              `json:"job_id"`
	VideoID     uint                `json:"video_id"`
	DownloadURL string              `json:"download_url"`
	Callback    VideoIngestCallback `json:"callback"`
}

type VideoIngestCallback struct {
	StatusURL        string `json:"status_url"`         // Webhook to report "processing", "ready", or "failed"
	UploadURL        string `json:"upload_url"`         // Pre-signed URL to PUT the master .m3u8 playlist
	RequestUploadURL string `json:"request_upload_url"` // Endpoint to request presigned URLs for individual segments (.ts)
}