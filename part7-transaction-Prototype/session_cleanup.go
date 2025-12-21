package main

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
    DeleteExpiredSessions() error
}

type DBSessionStore struct{}

func (s *DBSessionStore) DeleteExpiredSessions() error {
    // Implementation for deleting expired sessions from database
    log.Println("Deleting expired sessions from database")
    return nil
}

func cleanupSessions(store SessionStore) {
    if err := store.DeleteExpiredSessions(); err != nil {
        log.Printf("Failed to delete expired sessions: %v", err)
    }
}

func scheduleCleanup(store SessionStore, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        cleanupSessions(store)
    }
}

func main() {
    store := &DBSessionStore{}
    go scheduleCleanup(store, 24*time.Hour)

    // Keep main goroutine alive
    select {}
}