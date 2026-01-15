package main

import (
	"log"

	"github.com/Mahaveer86619/FrameSense/pkg/config"
	"github.com/Mahaveer86619/FrameSense/pkg/db"
	"github.com/Mahaveer86619/FrameSense/pkg/web"
)

func main() {
	config.LoadConfig()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	server := web.NewServer()
	web.StartServer(server)
}
