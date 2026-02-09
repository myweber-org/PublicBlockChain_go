
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("Cleaned up %d expired sessions", rowsAffected)
	return nil
}package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

const cleanupInterval = 1 * time.Hour
const sessionTTL = 24 * time.Hour

type SessionCleaner struct {
    db *sql.DB
}

func NewSessionCleaner(db *sql.DB) *SessionCleaner {
    return &SessionCleaner{db: db}
}

func (sc *SessionCleaner) Run(ctx context.Context) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            sc.cleanupExpiredSessions()
        }
    }
}

func (sc *SessionCleaner) cleanupExpiredSessions() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cutoffTime := time.Now().Add(-sessionTTL)

    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    result, err := sc.db.ExecContext(ctx, query, cutoffTime)
    if err != nil {
        log.Printf("Failed to clean expired sessions: %v", err)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Failed to get rows affected: %v", err)
        return
    }

    if rowsAffected > 0 {
        log.Printf("Cleaned up %d expired sessions", rowsAffected)
    }
}