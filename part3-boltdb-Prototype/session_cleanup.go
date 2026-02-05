
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
}package main

import (
	"context"
	"log"
	"time"

	"github.com/yourproject/internal/database"
)

const cleanupInterval = 24 * time.Hour

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	ctx := context.Background()
	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(ctx, db); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			} else {
				log.Println("Session cleanup completed successfully")
			}
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *database.DB) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rows)
	return nil
}