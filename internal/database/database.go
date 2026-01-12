package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Elmar006/api-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDB() {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "api_backend")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)

	log.Printf("Connecting to PostgreSQL: %s@%s:%s/%s", user, host, port, dbname)

	var db *gorm.DB
	var err error

	maxAttempts := 10
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			log.Printf("Successfully connected to PostgreSQL (attempt %d/%d)", i, maxAttempts)
			break
		}

		log.Printf("Attempt %d/%d failed: %v", i, maxAttempts, err)

		if i < maxAttempts {
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		log.Fatal("Failed to connect to database after all attempts:", err)
	}

	// Автомиграции
	log.Println("Running database migrations...")
	err = db.AutoMigrate(&models.Task{}, &models.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migrations completed")

	DB = Dbinstance{
		Db: db,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
