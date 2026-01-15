package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LocalStorageService struct {
	BasePath string
	BaseURL  string // Base URL for serving files (e.g., "http://localhost:7000/uploads")
}

func NewLocalStorageService() *LocalStorageService {
	return &LocalStorageService{
		BasePath: "uploads",
		BaseURL:  "http://localhost:8080/files", // todo configure via env
	}
}

func (l *LocalStorageService) ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

/* ---------- Save ---------- */

func (l *LocalStorageService) SaveVideo(
	file io.Reader,
	filename string,
) (string, error) {
	dir := filepath.Join(l.BasePath, "raw")
	if err := l.ensureDir(dir); err != nil {
		return "", err
	}
	return l.save(file, dir, filename)
}

func (l *LocalStorageService) SaveProcessedVideo(
	file io.Reader,
	filename string,
) (string, error) {
	dir := filepath.Join(l.BasePath, "processed")
	if err := l.ensureDir(dir); err != nil {
		return "", err
	}
	return l.save(file, dir, filename)
}

func (l *LocalStorageService) SaveHLSPlaylist(
	content io.Reader,
	filename string,
) (string, error) {
	dir := filepath.Join(l.BasePath, "hls")
	if err := l.ensureDir(dir); err != nil {
		return "", err
	}
	return l.save(content, dir, filename)
}

func (l *LocalStorageService) save(
	reader io.Reader,
	dir string,
	filename string,
) (string, error) {

	path := filepath.Join(dir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return path, nil
}

/* ---------- Get ---------- */

func (l *LocalStorageService) GetVideo(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (l *LocalStorageService) GetHLSPlaylist(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

/* ---------- Presigned URLs ---------- */

func (l *LocalStorageService) GeneratePresignedDownloadURL(
	path string,
	expiration time.Duration,
) (string, error) {
	// For local storage, return a simple URL
	// In production, implement token-based auth
	relativePath := filepath.ToSlash(path)
	return fmt.Sprintf("%s/%s", l.BaseURL, relativePath), nil
}

func (l *LocalStorageService) GeneratePresignedUploadURL(
	filename string,
	expiration time.Duration,
) (string, error) {
	// For local storage, return upload endpoint
	// In production, generate a signed token
	return fmt.Sprintf("%s/upload/processed/%s", l.BaseURL, filename), nil
}

/* ---------- Delete ---------- */

func (l *LocalStorageService) DeleteVideo(
	path string,
) error {
	return os.Remove(path)
}
