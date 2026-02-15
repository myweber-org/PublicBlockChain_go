
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

type SessionCleaner struct {
    db        *sql.DB
    interval  time.Duration
    retention time.Duration
}

func NewSessionCleaner(db *sql.DB, interval, retention time.Duration) *SessionCleaner {
    return &SessionCleaner{
        db:        db,
        interval:  interval,
        retention: retention,
    }
}

func (sc *SessionCleaner) Start(ctx context.Context) {
    ticker := time.NewTicker(sc.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            sc.cleanup()
        case <-ctx.Done():
            log.Println("Session cleaner stopped")
            return
        }
    }
}

func (sc *SessionCleaner) cleanup() {
    cutoff := time.Now().Add(-sc.retention)
    query := `DELETE FROM user_sessions WHERE last_activity < $1`

    result, err := sc.db.Exec(query, cutoff)
    if err != nil {
        log.Printf("Failed to clean sessions: %v", err)
        return
    }

    rows, _ := result.RowsAffected()
    log.Printf("Cleaned %d expired sessions", rows)
}