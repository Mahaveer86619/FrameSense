package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Mahaveer86619/FrameSense/pkg/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSQueueService struct {
	Client   *sqs.Client
	QueueURL string
}

func NewSQSQueueService(client *sqs.Client, queueURL string) *SQSQueueService {
	return &SQSQueueService{
		Client:   client,
		QueueURL: queueURL,
	}
}

func (s *SQSQueueService) SendVideoIngestMessage(msg *models.VideoIngestMessage) error {
	messageBody, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = s.Client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(string(messageBody)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"MessageType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("video.ingest"),
			},
			"VideoID": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprintf("%d", msg.VideoID)),
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	return nil
}
