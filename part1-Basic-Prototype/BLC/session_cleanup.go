package main

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
    DeleteExpiredSessions(ctx context.Context, before time.Time) (int64, error)
}

func cleanupExpiredSessions(ctx context.Context, store SessionStore, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Session cleanup stopped:", ctx.Err())
            return
        case <-ticker.C:
            before := time.Now()
            deleted, err := store.DeleteExpiredSessions(ctx, before)
            if err != nil {
                log.Printf("Failed to delete expired sessions: %v", err)
                continue
            }
            if deleted > 0 {
                log.Printf("Deleted %d expired sessions", deleted)
            }
        }
    }
}package main

import (
    "log"
    "time"
)

type Session struct {
    ID        string
    UserID    string
    ExpiresAt time.Time
}

type SessionStore interface {
    GetExpiredSessions() ([]Session, error)
    DeleteSession(id string) error
}

func cleanupExpiredSessions(store SessionStore) {
    expired, err := store.GetExpiredSessions()
    if err != nil {
        log.Printf("Failed to get expired sessions: %v", err)
        return
    }

    for _, session := range expired {
        if err := store.DeleteSession(session.ID); err != nil {
            log.Printf("Failed to delete session %s: %v", session.ID, err)
        } else {
            log.Printf("Deleted expired session: %s", session.ID)
        }
    }
}

func scheduleCleanup(store SessionStore, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        cleanupExpiredSessions(store)
    }
}

func main() {
    // Implementation would initialize SessionStore
    // and call scheduleCleanup with desired interval
    log.Println("Session cleanup scheduler started")
}