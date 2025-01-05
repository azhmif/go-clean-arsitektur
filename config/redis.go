package config

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),     // Contoh: "localhost:6379"
		Password: os.Getenv("REDIS_PASSWORD"), // Kosong jika tidak ada password
		DB:       0,                           // Database default
	})

	// Test koneksi Redis
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")
}
