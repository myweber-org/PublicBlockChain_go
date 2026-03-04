
package main

import (
    "context"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
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

    keys, err := rdb.ZRangeByScore(ctx, "user_sessions", &redis.ZRangeBy{
        Min: "-inf",
        Max: string(maxScore),
    }).Result()
    if err != nil {
        log.Printf("Failed to get expired sessions: %v", err)
        return
    }

    if len(keys) > 0 {
        pipe := rdb.Pipeline()
        for _, key := range keys {
            pipe.Del(ctx, "session:"+key)
        }
        pipe.ZRemRangeByScore(ctx, "user_sessions", "-inf", string(maxScore))
        _, err := pipe.Exec(ctx)
        if err != nil {
            log.Printf("Failed to delete expired sessions: %v", err)
        } else {
            log.Printf("Cleaned up %d expired sessions", len(keys))
        }
    }
}

func main() {
    initRedis()

    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions()
        }
    }
}