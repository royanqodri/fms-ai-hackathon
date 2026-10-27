package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/royanqodri/fms-ai-hackathon/config"
)

var RedisClient *redis.Client
var Ctx = context.Background()

const (
	KEY_API     = "api"
	KEY_SESSION = "session"
)

func InitRedis() {
	cfg := config.Get().Redis
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     20, // adjust as needed
		MinIdleConns: 5,  // keep some connections warm
	})

	// Retry until success (with backoff)
	for i := 0; i < 10; i++ { // max 10 retries
		_, err := RedisClient.Ping(Ctx).Result()
		if err == nil {
			log.Println("Connected to Redis successfully")
			return
		}
		log.Printf("Redis not ready, retrying in %d seconds...", i+1)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	log.Fatalf("Could not connect to Redis after retries")
}
