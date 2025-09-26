package util

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

const DEFAULT_HSET_GROUP = "gotask"

var (
	RedisData   *redis.Client
	RedisConfig *redis.Client
)

func CreateRedisConn(addr string, db int, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test Redis connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Warning: Could not connect to Redis: %v", err)
		log.Printf("Redis operations will fail. Please ensure Redis is running on localhost:6379")
		return nil, err
	} else {
		log.Printf("Successfully connected to Redis %s, db %d\n", addr, db)
	}
	return client, nil
}

func InitRedisClient(addr string, db int, dbConfig int, password string) error {
	client, err := CreateRedisConn(addr, db, password)
	if err != nil {
		return err
	} else {
		RedisData = client
	}

	client, err = CreateRedisConn(addr, dbConfig, password)
	if err != nil {
		return err
	} else {
		RedisConfig = client
	}
	return nil
}
