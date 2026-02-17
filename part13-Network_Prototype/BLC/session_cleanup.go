package main

import (
    "context"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    for {
        now := time.Now().Unix()
        maxScore := float64(now - 86400) // 24 hours ago

        // Remove expired sessions from sorted set
        removed, err := rdb.ZRemRangeByScore(ctx, "user_sessions", "0", 
            string(maxScore)).Result()
        if err != nil {
            log.Printf("Error cleaning sessions: %v", err)
        } else if removed > 0 {
            log.Printf("Cleaned %d expired sessions", removed)
        }

        // Wait 24 hours before next cleanup
        time.Sleep(24 * time.Hour)
    }
}