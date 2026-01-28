
package main

import (
    "context"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func initRedis() {
    rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
}

func cleanupExpiredSessions() {
    now := time.Now().Unix()
    maxScore := float64(now - 86400)

    // Remove sessions older than 24 hours
    removed, err := rdb.ZRemRangeByScore(ctx, "user_sessions", "0", string(maxScore)).Result()
    if err != nil {
        log.Printf("Failed to clean sessions: %v", err)
        return
    }

    log.Printf("Cleaned up %d expired sessions", removed)
}

func main() {
    initRedis()

    // Run cleanup daily at 2 AM
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions()
        }
    }
}