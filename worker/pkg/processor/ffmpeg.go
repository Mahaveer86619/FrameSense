package processor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Mahaveer86619/FrameSense-Worker/pkg/config"
	"github.com/Mahaveer86619/FrameSense-Worker/pkg/queue"
)

func ProcessVideo(job queue.VideoIngestMessage) error {
	cfg := config.AppConfig

	tmpDir := fmt.Sprintf("/tmp/%s", job.JobID)
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	inputFile := filepath.Join(tmpDir, "input.mp4")
	outputFile := filepath.Join(tmpDir, "output.m3u8")

	// 1. Download Video
	if err := downloadFile(job.DownloadURL, inputFile); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// 2. Select Encoder (CPU vs GPU)
	vCodec := "libx264"
	if cfg.UseGPU {
		switch cfg.GPUType {
		case "nvidia":
			vCodec = "h264_nvenc"
		case "vaapi", "intel":
			vCodec = "h264_vaapi"
		case "qsv":
			vCodec = "h264_qsv"
		default:
			vCodec = "h264_nvenc" // Default to NVIDIA
		}
	}

	// 3. Run FFmpeg (HLS Generation)
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-c:v", vCodec,
		"-preset", cfg.FFmpegPreset,
		"-g", "48", "-sc_threshold", "0",
		"-c:a", "aac", "-b:a", "128k",
		"-start_number", "0",
		"-hls_time", fmt.Sprintf("%d", cfg.HLSSegmentTime),
		"-hls_list_size", fmt.Sprintf("%d", cfg.HLSListSize),
		"-f", "hls", outputFile,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg failed: %s", string(output))
	}

	// 4. Upload Result (Using the pre-signed Upload URL for the playlist)
	return uploadFile(job.Callback.UploadURL, outputFile)
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func uploadFile(url, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("upload failed with status: %d", resp.StatusCode)
	}
	return nil
}

func SendCallback(url, status, errorMsg string) {
	// Payload matching server expectations
	body := []byte(fmt.Sprintf(`{"status":"%s", "error_message":"%s"}`, status, errorMsg))
	http.Post(url, "application/json", bytes.NewBuffer(body))
}
