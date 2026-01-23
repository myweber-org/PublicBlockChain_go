package main

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

	for {
		select {
		case <-ticker.C:
			cleanupExpiredSessions(db)
		}
	}
}

func cleanupExpiredSessions(db *database.Connection) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Failed to cleanup sessions: %v", err)
		return
	}

	rows, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rows)
}package main

import (
    "context"
    "log"
    "time"

    "github.com/yourproject/internal/database"
)

func main() {
    ctx := context.Background()
    db := database.NewConnection()

    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions(ctx, db)
        }
    }
}

func cleanupExpiredSessions(ctx context.Context, db *database.DB) {
    query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
    result, err := db.ExecContext(ctx, query)
    if err != nil {
        log.Printf("Failed to clean sessions: %v", err)
        return
    }

    rows, _ := result.RowsAffected()
    log.Printf("Cleaned %d expired sessions", rows)
}