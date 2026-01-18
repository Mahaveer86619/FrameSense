package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mahaveer86619/FrameSense-Worker/pkg/config"
	"github.com/Mahaveer86619/FrameSense-Worker/pkg/processor"
	"github.com/Mahaveer86619/FrameSense-Worker/pkg/queue"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	cfg.PrintConfig()

	// Create job queue channel
	jobQueue := make(chan queue.VideoIngestMessage, 100)

	// Initialize SQS Client
	sqsClient := queue.NewSQSClient()

	// Start Worker Pool
	pool := processor.NewWorkerPool(jobQueue)
	pool.Start()

	// Start SQS Polling in a separate goroutine
	ctx, cancel := context.WithCancel(context.Background())
	go queue.PollSQS(ctx, sqsClient, cfg.SQSQueueURL, jobQueue)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, stopping worker...")

	cancel()        // Stop SQS polling
	pool.Stop()     // Stop worker pool
	close(jobQueue) // Close job queue

	log.Println("Worker service stopped gracefully")
}
