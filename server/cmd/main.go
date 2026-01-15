package main

import (
	"github.com/Mahaveer86619/FrameSense/pkg/config"
	"github.com/Mahaveer86619/FrameSense/pkg/web"
)

func main() {
	config.LoadConfig()

	server := web.NewServer()
	web.StartServer(server)
}
