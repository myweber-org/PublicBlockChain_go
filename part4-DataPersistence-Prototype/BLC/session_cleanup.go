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
}