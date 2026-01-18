package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Mahaveer86619/FrameSense-Worker/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewSQSClient() *sqs.Client {
	cfg := config.AppConfig

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKey,
			cfg.AWSSecretKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(awsCfg, func(o *sqs.Options) {
		if cfg.SQSEndpoint != "" {
			o.BaseEndpoint = aws.String(cfg.SQSEndpoint)
			log.Printf("Using SQS endpoint: %s", cfg.SQSEndpoint)
		}
	})

	return client
}

func PollSQS(ctx context.Context, client *sqs.Client, queueURL string, jobQueue chan<- VideoIngestMessage) {
	log.Printf("Started polling SQS queue: %s", queueURL)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping SQS polling...")
			return
		default:
			// 1. Receive Message
			output, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueURL),
				MaxNumberOfMessages: 1,   // Fetch one job at a time per poll
				WaitTimeSeconds:     20,  // Long Polling (reduces empty responses)
				VisibilityTimeout:   300, // 5 minutes to process before it becomes visible again
			})

			if err != nil {
				log.Printf("Error fetching messages: %v", err)
				time.Sleep(5 * time.Second) // Backoff on error
				continue
			}

			if len(output.Messages) == 0 {
				continue
			}

			// 2. Process Message
			for _, msg := range output.Messages {
				var ingestMsg VideoIngestMessage
				if err := json.Unmarshal([]byte(*msg.Body), &ingestMsg); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					// Strictly delete malformed messages to prevent infinite loops
					deleteMessage(client, queueURL, msg.ReceiptHandle)
					continue
				}

				log.Printf("Received Job: %s (VideoID: %d)", ingestMsg.JobID, ingestMsg.VideoID)

				// 3. Send to Worker Pool
				jobQueue <- ingestMsg

				// 4. Delete from SQS (Ack)
				deleteMessage(client, queueURL, msg.ReceiptHandle)
			}
		}
	}
}

func deleteMessage(client *sqs.Client, queueURL string, receiptHandle *string) {
	_, err := client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
	}
}
