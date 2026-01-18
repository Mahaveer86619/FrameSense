package processor

import (
	"log"
	"sync"

	"github.com/Mahaveer86619/FrameSense-Worker/pkg/config"
	"github.com/Mahaveer86619/FrameSense-Worker/pkg/queue"
)

type WorkerPool struct {
	Workers  int
	JobQueue <-chan queue.VideoIngestMessage
	wg       sync.WaitGroup
	quit     chan bool
}

func NewWorkerPool(jobQueue <-chan queue.VideoIngestMessage) *WorkerPool {
	cfg := config.AppConfig

	return &WorkerPool{
		Workers:  cfg.WorkerCount,
		JobQueue: jobQueue,
		quit:     make(chan bool),
	}
}

func (wp *WorkerPool) Start() {
	cfg := config.AppConfig

	log.Printf("Worker service started with %d workers. GPU Enabled: %t",
		cfg.WorkerCount, cfg.UseGPU)

	if cfg.UseGPU {
		log.Printf("Using GPU acceleration: %s", cfg.GPUType)
	}

	for i := 0; i < wp.Workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	log.Printf("[Worker %d] Started", id)

	for {
		select {
		case job, ok := <-wp.JobQueue:
			if !ok {
				return
			}
			log.Printf("[Worker %d] Processing Job: %s", id, job.JobID)

			// Process Video (Implemented in ffmpeg.go)
			err := ProcessVideo(job)

			// Report Status back to Main Server
			if err != nil {
				log.Printf("[Worker %d] Failed: %v", id, err)
				SendCallback(job.Callback.StatusURL, "failed", err.Error())
			} else {
				log.Printf("[Worker %d] Success", id)
				SendCallback(job.Callback.StatusURL, "ready", "")
			}

		case <-wp.quit:
			return
		}
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.quit)
	wp.wg.Wait()
}
