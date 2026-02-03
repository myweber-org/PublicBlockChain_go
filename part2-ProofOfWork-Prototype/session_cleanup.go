
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
	"yourproject/internal/models"
)

const cleanupInterval = 24 * time.Hour

func main() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		cleanupExpiredSessions()
	}
}

func cleanupExpiredSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.Conn.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Failed to clean up sessions: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Cleaned up %d expired sessions", rowsAffected)
	}
}package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
	"yourproject/internal/models"
)

func main() {
	ctx := context.Background()
	database := db.GetDB()

	// Run cleanup every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupExpiredSessions(ctx, database)
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *db.Database) {
	cutoff := time.Now().Add(-24 * time.Hour)
	result := db.WithContext(ctx).
		Where("last_activity < ?", cutoff).
		Delete(&models.Session{})

	if result.Error != nil {
		log.Printf("Error cleaning sessions: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d expired sessions", result.RowsAffected)
	}
}package main

import (
    "context"
    "log"
    "time"

    "github.com/jackc/pgx/v5"
)

const (
    cleanupInterval = 1 * time.Hour
    sessionTTL      = 24 * time.Hour
    deleteBatchSize = 1000
)

func main() {
    connStr := "postgresql://user:pass@localhost:5432/dbname"
    conn, err := pgx.Connect(context.Background(), connStr)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }
    defer conn.Close(context.Background())

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for range ticker.C {
        if err := cleanupExpiredSessions(conn); err != nil {
            log.Printf("Session cleanup failed: %v", err)
        }
    }
}

func cleanupExpiredSessions(conn *pgx.Conn) error {
    ctx := context.Background()
    cutoffTime := time.Now().Add(-sessionTTL)

    for {
        tag, err := conn.Exec(ctx,
            `DELETE FROM user_sessions 
             WHERE last_activity < $1 
             AND session_id IN (
                 SELECT session_id FROM user_sessions 
                 WHERE last_activity < $1 
                 LIMIT $2
             )`,
            cutoffTime, deleteBatchSize)

        if err != nil {
            return err
        }

        if tag.RowsAffected() == 0 {
            break
        }

        log.Printf("Deleted %d expired sessions", tag.RowsAffected())
        time.Sleep(100 * time.Millisecond)
    }

    return nil
}