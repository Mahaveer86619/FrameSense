package main

import (
	"log"

	"github.com/Mahaveer86619/FrameSense/pkg/config"
	"github.com/Mahaveer86619/FrameSense/pkg/db"
	"github.com/Mahaveer86619/FrameSense/pkg/models"
)

func main() {
	config.LoadConfig()

	err := db.InitDB(true)
	if err != nil {
		log.Panicf("Error while db Init")
	}

	err = db.MigrateDB(
		&models.User{},
		&models.Video{},
	)
	if err != nil {
		log.Panicf("Error while db migration")
	}
}
