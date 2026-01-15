package consts

type VideoStatus string

const (
	VideoStatusUploaded    VideoStatus = "uploaded"
	VideoStatusQueued      VideoStatus = "queued"
	VideoStatusProcessing  VideoStatus = "processing"
	VideoStatusReady       VideoStatus = "ready"
	VideoStatusQueueFailed VideoStatus = "queue_failed"
	VideoStatusFailed      VideoStatus = "failed"
)
