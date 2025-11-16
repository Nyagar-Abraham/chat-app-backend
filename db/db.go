package db

import (
	"fmt"
	"log"
	"os"

	"github.com/Nyagar-Abraham/chat-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// IsRecordNotFoundError checks if the error is a "record not found" error
func IsRecordNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	// Conditionally run AutoMigrate if MIGRATE_DB=true in env (for development only)
	if os.Getenv("MIGRATE_DB") == "true" {
		fmt.Println("[DEV] Running GORM AutoMigration...")
		err = db.AutoMigrate(&models.Tenant{}, &models.Channel{}, &models.User{}, &models.ChannelMember{})
		if err != nil {
			log.Fatalf("Failed to run migration: %v", err)
		}
		fmt.Println("Database connected and migrated (AutoMigrate enabled)")
	} else {
		fmt.Println("Database connected (no migration performed)")
	}
}
