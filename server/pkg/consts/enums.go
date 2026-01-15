package consts

type VideoStatus string

const (
	VideoStatusUploaded   VideoStatus = "uploaded"
	VideoStatusProcessing VideoStatus = "processing"
	VideoStatusReady      VideoStatus = "ready"
	VideoStatusFailed     VideoStatus = "failed"
)
