package database

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/storage/redis/v3"
	rds "github.com/redis/go-redis/v9"
)

func NewRedisStorage() *redis.Storage {
	storage := redis.New(redis.Config{
		Host:     "0.0.0.0",
		Port:     6379,
		Password: "",
		Username: "",
		Database: 0,
	})
	return storage
}

var redisClient *rds.Client

// InitRedis initializes the Redis client connection
func InitRedis(addr, password string, db int) error {
	redisClient = rds.NewClient(&rds.Options{
		Addr:     "0.0.0.0:6379",
		Password: "",
		DB:       0,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return err
	}

	log.Println("Successfully connected to Redis")
	return nil
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *rds.Client {
	return redisClient
}

// CloseRedis closes the Redis client connection
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
