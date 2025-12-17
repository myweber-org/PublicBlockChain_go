package main

import (
    "context"
    "log"
    "time"
)

const (
    cleanupInterval = 1 * time.Hour
    sessionTTL      = 24 * time.Hour
)

type SessionStore interface {
    DeleteExpiredSessions(ctx context.Context, olderThan time.Time) (int64, error)
}

func startSessionCleanup(store SessionStore, stopCh <-chan struct{}) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupSessions(store)
        case <-stopCh:
            log.Println("Session cleanup stopped")
            return
        }
    }
}

func cleanupSessions(store SessionStore) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cutoff := time.Now().Add(-sessionTTL)
    deletedCount, err := store.DeleteExpiredSessions(ctx, cutoff)
    if err != nil {
        log.Printf("Failed to clean up sessions: %v", err)
        return
    }

    if deletedCount > 0 {
        log.Printf("Cleaned up %d expired sessions", deletedCount)
    }
}