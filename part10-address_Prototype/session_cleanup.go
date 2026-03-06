
package main

import (
    "context"
    "log"
    "time"

    "github.com/yourproject/db"
)

const cleanupInterval = 24 * time.Hour
const sessionTTL = 7 * 24 * time.Hour

func cleanupExpiredSessions() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    cutoff := time.Now().Add(-sessionTTL)
    result, err := db.Conn.ExecContext(ctx,
        "DELETE FROM user_sessions WHERE last_activity < ?", cutoff)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rows)
    return nil
}

func startCleanupScheduler() {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for range ticker.C {
        if err := cleanupExpiredSessions(); err != nil {
            log.Printf("Session cleanup failed: %v", err)
        }
    }
}

func main() {
    if err := db.Initialize(); err != nil {
        log.Fatal("Database initialization failed:", err)
    }
    defer db.Close()

    go startCleanupScheduler()
    select {}
}package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

const (
    cleanupInterval = 1 * time.Hour
    sessionTTL      = 24 * time.Hour
)

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

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected > 0 {
        log.Printf("Cleaned %d expired sessions", rowsAffected)
    }
}