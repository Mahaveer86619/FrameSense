package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorageService struct {
	BasePath string
}

func NewLocalStorageService() *LocalStorageService {
	return &LocalStorageService{
		BasePath: "uploads",
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

/* ---------- Delete ---------- */

func (l *LocalStorageService) DeleteVideo(
	path string,
) error {
	return os.Remove(path)
}
