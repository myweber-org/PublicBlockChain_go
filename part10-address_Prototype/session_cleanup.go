
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
}