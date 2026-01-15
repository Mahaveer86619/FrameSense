package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	c "github.com/Mahaveer86619/FrameSense/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageService interface {
	SaveVideo(file io.Reader, filename string) (string, error)

	GetVideo(path string) (io.ReadCloser, error)
	GetHLSPlaylist(path string) (io.ReadCloser, error)

	SaveProcessedVideo(file io.Reader, filename string) (string, error)
	SaveHLSPlaylist(content io.Reader, filename string) (string, error)

	GeneratePresignedDownloadURL(path string, expiration time.Duration) (string, error)
	GeneratePresignedUploadURL(filename string, expiration time.Duration) (string, error)

	DeleteVideo(path string) error
}

func NewStorageService() (StorageService, error) {
	switch c.AppConfig.STORAGE_DRIVER {
	case "local":
		return NewLocalStorageService(), nil

	case "s3":
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.AppConfig.S3_REGION),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				c.AppConfig.S3_ACCESS_KEY,
				c.AppConfig.S3_SECRET_KEY,
				"",
			)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load aws config: %w", err)
		}

		client := s3.NewFromConfig(cfg, func(o *s3.Options) {
			if c.AppConfig.S3_ENDPOINT != "" {
				o.BaseEndpoint = aws.String(c.AppConfig.S3_ENDPOINT)
			}
			o.UsePathStyle = true
		})

		bucketName := c.AppConfig.S3_BUCKET
		_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})

		if err != nil {
			_, err = client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
				Bucket: aws.String(bucketName),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to ensure bucket %s exists: %w", bucketName, err)
			}
		}

		return NewS3StorageService(client, c.AppConfig.S3_BUCKET), nil

	default:
		return nil, fmt.Errorf("unsupported storage type: %s", c.AppConfig.STORAGE_DRIVER)
	}
}
