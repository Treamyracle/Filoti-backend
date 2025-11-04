package config

import (
	"fmt"
	"log"
	"os"

	"filoti-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST_POOLER")
	port := os.Getenv("DB_PORT_POOLER")
	dbname := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("One or more required database environment variables (DB_USER, DB_PASSWORD, DB_HOST_POOLER, DB_PORT_POOLER, DB_NAME) are not set. Please check your .env file.")
	}

	sslmode := os.Getenv("DB_SSLMODE_POOLER")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta&pgbouncer=true",
		host, user, password, dbname, port, sslmode,
	)

	log.Printf("Attempting to connect to database using DSN (Pooler): %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database!")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Status{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Database connected and migrated successfully.")
}
