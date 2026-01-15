package queue

import (
	"context"
	"fmt"

	c "github.com/Mahaveer86619/FrameSense/pkg/config"
	"github.com/Mahaveer86619/FrameSense/pkg/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type QueueService interface {
	SendVideoIngestMessage(msg *models.VideoIngestMessage) error
	// Add more message types as needed
	// SendAdGenerationMessage(msg *AdGenerationMessage) error
	// SendSceneAnalysisMessage(msg *SceneAnalysisMessage) error
}

func NewQueueService() (QueueService, error) {
	switch c.AppConfig.QUEUE_DRIVER {
	case "sqs":
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.AppConfig.SQS_REGION),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				c.AppConfig.SQS_ACCESS_KEY,
				c.AppConfig.SQS_SECRET_KEY,
				"",
			)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load aws config: %w", err)
		}

		client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			if c.AppConfig.SQS_ENDPOINT != "" {
				o.BaseEndpoint = aws.String(c.AppConfig.SQS_ENDPOINT)
			}
		})

		queueName := c.AppConfig.SQS_QUEUE_NAME
		queueURL, err := ensureQueueExists(client, queueName)
		if err != nil {
			return nil, fmt.Errorf("failed to ensure queue exists: %w", err)
		}

		return NewSQSQueueService(client, queueURL), nil

	case "local":
		return NewLocalQueueService(), nil

	default:
		return nil, fmt.Errorf("unsupported queue driver: %s", c.AppConfig.QUEUE_DRIVER)
	}
}

func ensureQueueExists(client *sqs.Client, queueName string) (string, error) {
	getQueueURLOutput, err := client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	if err == nil {
		return *getQueueURLOutput.QueueUrl, nil
	}

	createQueueOutput, err := client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]string{
			"VisibilityTimeout":      "300",
			"MessageRetentionPeriod": "86400",
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create queue: %w", err)
	}

	return *createQueueOutput.QueueUrl, nil
}
