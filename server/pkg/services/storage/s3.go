package storage

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3StorageService struct {
	Client         *s3.Client
	Bucket         string
	PresignClient  *s3.PresignClient
}

func NewS3StorageService(client *s3.Client, bucket string) *S3StorageService {
	return &S3StorageService{
		Client:        client,
		Bucket:        bucket,
		PresignClient: s3.NewPresignClient(client),
	}
}

/* ---------- Save ---------- */

func (s *S3StorageService) SaveVideo(
	file io.Reader,
	filename string,
) (string, error) {
	return s.put("raw/"+filename, file)
}

func (s *S3StorageService) SaveProcessedVideo(
	file io.Reader,
	filename string,
) (string, error) {
	return s.put("processed/"+filename, file)
}

func (s *S3StorageService) SaveHLSPlaylist(
	content io.Reader,
	filename string,
) (string, error) {
	return s.put("hls/"+filename, content)
}

func (s *S3StorageService) put(
	key string,
	body io.Reader,
) (string, error) {

	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
		Body:   body,
	})
	if err != nil {
		return "", err
	}

	return "s3://" + s.Bucket + "/" + key, nil
}

/* ---------- Get ---------- */

func (s *S3StorageService) GetVideo(
	path string,
) (io.ReadCloser, error) {
	return s.get(path)
}

func (s *S3StorageService) GetHLSPlaylist(
	path string,
) (io.ReadCloser, error) {
	return s.get(path)
}

func (s *S3StorageService) get(
	path string,
) (io.ReadCloser, error) {

	key := strings.TrimPrefix(path, "s3://"+s.Bucket+"/")

	resp, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

/* ---------- Generate Presigned URLs ---------- */

func (s *S3StorageService) GeneratePresignedDownloadURL(
	path string,
	expiration time.Duration,
) (string, error) {
	key := strings.TrimPrefix(path, "s3://"+s.Bucket+"/")

	req, err := s.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (s *S3StorageService) GeneratePresignedUploadURL(
	filename string,
	expiration time.Duration,
) (string, error) {
	key := "processed/" + filename

	req, err := s.PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

/* ---------- Delete ---------- */

func (s *S3StorageService) DeleteVideo(
	path string,
) error {

	// path = s3://bucket/key
	key := strings.TrimPrefix(path, "s3://"+s.Bucket+"/")

	_, err := s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	return err
}
