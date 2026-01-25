package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func cleanupExpiredSessions() error {
	now := time.Now().Unix()
	keys, err := rdb.Keys(ctx, "session:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		exp, err := rdb.Get(ctx, key+":expires").Int64()
		if err != nil {
			continue
		}
		if exp < now {
			rdb.Del(ctx, key, key+":expires")
			log.Printf("Removed expired session: %s", key)
		}
	}
	return nil
}

func main() {
	initRedis()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(); err != nil {
				log.Printf("Cleanup failed: %v", err)
			}
		}
	}
}