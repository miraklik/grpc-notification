package db

import (
	"fmt"
	"log"
	"notification_service/internal/config"
	"notification_service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}

	if err := db.AutoMigrate(&models.Notification{}); err != nil {
		log.Printf("Error migrating database: %v", err)
		return nil, err
	}

	return db, nil
}
