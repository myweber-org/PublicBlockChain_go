package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

const (
    sessionTTL = 7 * 24 * time.Hour
    batchSize  = 100
)

func cleanupExpiredSessions(ctx context.Context, db *sql.DB) error {
    cutoff := time.Now().Add(-sessionTTL)
    query := `DELETE FROM user_sessions WHERE last_activity < $1 LIMIT $2`

    for {
        result, err := db.ExecContext(ctx, query, cutoff, batchSize)
        if err != nil {
            return err
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            return err
        }

        if rowsAffected == 0 {
            break
        }

        log.Printf("Deleted %d expired sessions", rowsAffected)

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(100 * time.Millisecond):
        }
    }

    return nil
}

func main() {
    db, err := sql.Open("postgres", "postgresql://localhost/sessions")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx := context.Background()
    if err := cleanupExpiredSessions(ctx, db); err != nil {
        log.Fatal(err)
    }

    log.Println("Session cleanup completed")
}package main

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

    rows, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rows)
    return nil
}