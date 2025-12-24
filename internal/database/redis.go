package database

import (
	"context"
	"log"
	"main/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RDB *redis.Client
func InitRedis(url string) {

	ctx:=context.Background()
    opt, err := redis.ParseURL(url)
    if err != nil {
        log.Fatalf("Failed to parse Redis URL: %v", err)
    }

	RDB=redis.NewClient(opt)
    pong, err := RDB.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }
	logger.Log.Info("Redis client connected ",zap.String("Ping response",pong))
}