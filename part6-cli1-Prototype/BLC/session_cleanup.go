package main

import (
    "context"
    "database/sql"
    "log"
    "time"

    _ "github.com/lib/pq"
)

const (
    dbConnectionString = "postgres://user:password@localhost/sessions?sslmode=disable"
    cleanupInterval    = 1 * time.Hour
    sessionTTL         = 24 * time.Hour
)

func main() {
    db, err := sql.Open("postgres", dbConnectionString)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    ctx := context.Background()

    for {
        select {
        case <-ticker.C:
            err := cleanupExpiredSessions(ctx, db)
            if err != nil {
                log.Printf("Session cleanup failed: %v", err)
            } else {
                log.Println("Session cleanup completed successfully")
            }
        }
    }
}

func cleanupExpiredSessions(ctx context.Context, db *sql.DB) error {
    cutoffTime := time.Now().Add(-sessionTTL)

    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    result, err := db.ExecContext(ctx, query, cutoffTime)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    log.Printf("Removed %d expired sessions", rowsAffected)
    return nil
}