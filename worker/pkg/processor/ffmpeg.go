package processor

import (
	"bytes"
	"encoding/json"
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
	outputDir := filepath.Join(tmpDir, "hls")
	os.MkdirAll(outputDir, 0755)

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
			vCodec = "h264_nvenc"
		}
	}

	// 3. Generate HLS with multiple quality variants
	// This creates a master playlist and multiple quality variants
	masterPlaylist := filepath.Join(outputDir, "master.m3u8")

	// Generate multiple quality variants
	variants := []struct {
		name         string
		resolution   string
		videoBitrate string
		audioBitrate string
	}{
		{"720p", "1280x720", "2800k", "128k"},
		{"480p", "854x480", "1400k", "128k"},
		{"360p", "640x360", "800k", "96k"},
	}

	// Build FFmpeg command with multiple outputs
	args := []string{
		"-i", inputFile,
		"-c:v", vCodec,
		"-preset", cfg.FFmpegPreset,
		"-g", "48",
		"-sc_threshold", "0",
		"-c:a", "aac",
	}

	// Add variant streams
	var variantPlaylists []string
	for i, variant := range variants {
		variantPlaylist := filepath.Join(outputDir, fmt.Sprintf("%s.m3u8", variant.name))
		variantPlaylists = append(variantPlaylists, variant.name)

		args = append(args,
			// Map the input video and audio
			"-map", "0:v:0",
			"-map", "0:a:0",
			// Set video filters for resolution
			fmt.Sprintf("-s:v:%d", i), variant.resolution,
			// Set bitrates
			fmt.Sprintf("-b:v:%d", i), variant.videoBitrate,
			fmt.Sprintf("-b:a:%d", i), variant.audioBitrate,
			// HLS specific settings
			"-hls_time", fmt.Sprintf("%d", cfg.HLSSegmentTime),
			"-hls_list_size", fmt.Sprintf("%d", cfg.HLSListSize),
			"-hls_segment_filename", filepath.Join(outputDir, fmt.Sprintf("%s_%%03d.ts", variant.name)),
			variantPlaylist,
		)
	}

	// Run FFmpeg
	cmd := exec.Command("ffmpeg", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg failed: %s", string(output))
	}

	// 4. Create master playlist
	if err := createMasterPlaylist(masterPlaylist, variants, outputDir); err != nil {
		return fmt.Errorf("failed to create master playlist: %w", err)
	}

	// 5. Upload all HLS files to the server
	if err := uploadHLSFiles(job, outputDir); err != nil {
		return fmt.Errorf("failed to upload HLS files: %w", err)
	}

	return nil
}

func createMasterPlaylist(
	path string,
	variants []struct {
		name         string
		resolution   string
		videoBitrate string
		audioBitrate string
	}, outputDir string,
) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write master playlist header
	f.WriteString("#EXTM3U\n")
	f.WriteString("#EXT-X-VERSION:3\n\n")

	// Write variant streams
	for _, variant := range variants {
		// Parse bitrate (remove 'k' suffix and convert to bits)
		var bitrate int
		fmt.Sscanf(variant.videoBitrate, "%dk", &bitrate)
		bitrate *= 1000 // Convert to bits

		f.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n", bitrate, variant.resolution))
		f.WriteString(fmt.Sprintf("%s.m3u8\n\n", variant.name))
	}

	return nil
}

func uploadHLSFiles(job queue.VideoIngestMessage, hlsDir string) error {
	// Walk through all files in the HLS directory
	return filepath.Walk(hlsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path from hlsDir
		relPath, err := filepath.Rel(hlsDir, path)
		if err != nil {
			return err
		}

		// Request upload URL for this file
		uploadURL, err := requestUploadURL(job.Callback.RequestUploadURL, relPath)
		if err != nil {
			return fmt.Errorf("failed to get upload URL for %s: %w", relPath, err)
		}

		// Upload the file
		if err := uploadFile(uploadURL, path); err != nil {
			return fmt.Errorf("failed to upload %s: %w", relPath, err)
		}

		return nil
	})
}

func requestUploadURL(baseURL, filename string) (string, error) {
	url := fmt.Sprintf("%s?filename=%s", baseURL, filename)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get upload URL: status %d", resp.StatusCode)
	}

	var result struct {
		UploadURL string `json:"upload_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.UploadURL, nil
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
	body := []byte(fmt.Sprintf(`{"status":"%s", "error_message":"%s"}`, status, errorMsg))
	http.Post(url, "application/json", bytes.NewBuffer(body))
}
