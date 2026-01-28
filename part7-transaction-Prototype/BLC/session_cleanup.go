package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

const cleanupInterval = 24 * time.Hour
const sessionTTL = 7 * 24 * time.Hour

func cleanupExpiredSessions(db *sql.DB) {
    for {
        time.Sleep(cleanupInterval)
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
        defer cancel()

        query := `DELETE FROM user_sessions WHERE last_activity < $1`
        cutoffTime := time.Now().Add(-sessionTTL)
        result, err := db.ExecContext(ctx, query, cutoffTime)
        if err != nil {
            log.Printf("Failed to clean up sessions: %v", err)
            continue
        }

        rowsAffected, _ := result.RowsAffected()
        log.Printf("Cleaned up %d expired sessions", rowsAffected)
    }
}

func main() {
    db, err := sql.Open("postgres", "postgresql://localhost/sessions?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    cleanupExpiredSessions(db)
}