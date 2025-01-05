package config

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection
func InitDB() *gorm.DB {
	// Database configuration
	dsn := getDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}

	log.Println("Database connection established")
	return db
}

// getDSN constructs the Data Source Name (DSN) for MySQL
func getDSN() string {
	// Load environment variables or provide default values
	username := os.Getenv("DB_USER")
	if username == "" {
		username = "root"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = ""
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}

	database := os.Getenv("DB_NAME")
	if database == "" {
		database = "crud_db"
	}

	return username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
}
