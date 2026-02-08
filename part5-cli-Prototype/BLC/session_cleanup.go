package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

const cleanupInterval = 1 * time.Hour
const sessionTTL = 24 * time.Hour

func cleanupExpiredSessions(db *sql.DB) {
    ctx := context.Background()
    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    cutoff := time.Now().Add(-sessionTTL)

    result, err := db.ExecContext(ctx, query, cutoff)
    if err != nil {
        log.Printf("Failed to clean up sessions: %v", err)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
}

func startSessionCleanupJob(db *sql.DB) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions(db)
        }
    }
}

func main() {
    db, err := sql.Open("postgres", "postgresql://localhost/sessions")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    startSessionCleanupJob(db)
}package main

import (
    "log"
    "time"
    "database/sql"
    _ "github.com/lib/pq"
)

func cleanupExpiredSessions(db *sql.DB) error {
    query := `DELETE FROM user_sessions WHERE expires_at < $1`
    result, err := db.Exec(query, time.Now())
    if err != nil {
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
    return nil
}

func scheduleCleanup(db *sql.DB) {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := cleanupExpiredSessions(db); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            }
        }
    }
}

func main() {
    db, err := sql.Open("postgres", "host=localhost user=app dbname=appdb sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    go scheduleCleanup(db)
    
    select {}
}