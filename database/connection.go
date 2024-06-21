package database

import (
	"log"
	"server/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

func ConnectDB() {
	connectionString := config.Env("DB_CONNECTION_STRING")

	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	log.Println("Database connected")

	// Create pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error creating pool: %v\n", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}
