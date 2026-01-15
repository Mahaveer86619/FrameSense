package queue

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/Mahaveer86619/FrameSense/pkg/models"
)

type LocalQueueService struct {
	messages []interface{}
	mu       sync.RWMutex
}

func NewLocalQueueService() *LocalQueueService {
	return &LocalQueueService{
		messages: make([]interface{}, 0),
	}
}

func (l *LocalQueueService) SendVideoIngestMessage(msg *models.VideoIngestMessage) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.messages = append(l.messages, msg)

	// Log for debugging/testing purposes
	msgJSON, _ := json.MarshalIndent(msg, "", "  ")
	log.Printf("[LOCAL QUEUE] Video Ingest Message Queued:\n%s\n", string(msgJSON))

	return nil
}

// GetMessages returns all queued messages (useful for testing)
func (l *LocalQueueService) GetMessages() []interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	messages := make([]interface{}, len(l.messages))
	copy(messages, l.messages)
	return messages
}

// ClearMessages clears all messages (useful for testing)
func (l *LocalQueueService) ClearMessages() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.messages = make([]interface{}, 0)
}
