package main

import (
    "log"
    "time"
    "context"
    "database/sql"
    _ "github.com/lib/pq"
)

func cleanupExpiredSessions(db *sql.DB) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    query := `DELETE FROM user_sessions WHERE expires_at < $1`
    result, err := db.ExecContext(ctx, query, time.Now())
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    log.Printf("Cleaned up %d expired sessions", rowsAffected)
    return nil
}

func main() {
    db, err := sql.Open("postgres", "postgres://user:pass@localhost/dbname?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

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
}package main

import (
    "context"
    "log"
    "time"
)

type Session struct {
    ID        string
    UserID    string
    ExpiresAt time.Time
}

type SessionStore interface {
    DeleteExpiredSessions(ctx context.Context) error
}

func cleanupExpiredSessions(store SessionStore, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        err := store.DeleteExpiredSessions(ctx)
        cancel()

        if err != nil {
            log.Printf("Failed to delete expired sessions: %v", err)
        } else {
            log.Println("Successfully cleaned up expired sessions")
        }
    }
}

func main() {
    // In a real application, initialize your session store here
    // store := NewDatabaseSessionStore()
    // cleanupExpiredSessions(store, 1*time.Hour)
    
    log.Println("Session cleanup service would start here")
    select {} // Block forever in this example
}