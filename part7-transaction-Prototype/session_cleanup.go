package main

import (
    "context"
    "log"
    "time"

    "yourproject/internal/database"
    "yourproject/internal/models"
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

    log.Println("Session cleanup service started")

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions(db)
        }
    }
}

func cleanupExpiredSessions(db *database.DB) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    result, err := db.ExecContext(ctx,
        "DELETE FROM user_sessions WHERE expires_at < NOW()",
    )
    if err != nil {
        log.Printf("Failed to cleanup sessions: %v", err)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
}