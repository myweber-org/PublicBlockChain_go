package main

import (
    "context"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
)

const (
    sessionPrefix = "session:"
    batchSize     = 100
)

func cleanupExpiredSessions(ctx context.Context, client *redis.Client) error {
    var cursor uint64
    var keys []string
    var err error

    for {
        keys, cursor, err = client.Scan(ctx, cursor, sessionPrefix+"*", batchSize).Result()
        if err != nil {
            return err
        }

        if len(keys) > 0 {
            pipe := client.Pipeline()
            for _, key := range keys {
                pipe.Exists(ctx, key)
            }
            cmds, err := pipe.Exec(ctx)
            if err != nil {
                log.Printf("Pipeline execution failed: %v", err)
                continue
            }

            var toDelete []string
            for i, cmd := range cmds {
                exists := cmd.(*redis.IntCmd).Val()
                if exists == 0 {
                    toDelete = append(toDelete, keys[i])
                }
            }

            if len(toDelete) > 0 {
                if err := client.Del(ctx, toDelete...).Err(); err != nil {
                    log.Printf("Failed to delete keys: %v", err)
                } else {
                    log.Printf("Deleted %d expired sessions", len(toDelete))
                }
            }
        }

        if cursor == 0 {
            break
        }
    }
    return nil
}

func main() {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    defer rdb.Close()

    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := cleanupExpiredSessions(ctx, rdb); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            }
        }
    }
}