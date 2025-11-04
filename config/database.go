package config

import (
	"log"
	"os"

	"filoti-backend/models" // Pastikan path ini sesuai dengan struktur proyek Anda

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	// Muat file .env, abaikan jika tidak ada (misalnya dalam produksi)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Ambil DATABASE_URL dari environment
	// Neon menyediakan ini sebagai satu connection string lengkap
	dsn := os.Getenv("DATABASE_URL")

	// Validasi bahwa DATABASE_URL ada
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set. Please check your .env file or environment configuration.")
	}

	log.Println("Attempting to connect to database using DATABASE_URL...")

	// Buka koneksi GORM menggunakan DSN dari DATABASE_URL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Tetapkan instance DB global
	DB = db

	// Dapatkan objek sql.DB dasar untuk melakukan Ping
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}

	// Lakukan Ping untuk memastikan koneksi benar-benar hidup
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database!")

	// Jalankan AutoMigrate seperti sebelumnya
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
