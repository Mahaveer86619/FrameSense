package db

import (
	"fmt"
	"log"

	"github.com/Mahaveer86619/FrameSense/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(local ...bool) error {
	var host string
	if len(local) > 0 && local[0] {
		host = "localhost"
	} else {
		host = config.AppConfig.DB_HOST
	}

	dsn := "host=" + host +
		" user=" + config.AppConfig.DB_USER +
		" password=" + config.AppConfig.DB_PASSWORD +
		" dbname=" + config.AppConfig.DB_NAME +
		" port=" + fmt.Sprintf("%s", config.AppConfig.DB_PORT) +
		" sslmode=disable TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return err
}

func GetDB() *gorm.DB {
	return DB
}

func MigrateDB(models ...interface{}) error {
	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Printf("Error during migration: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
