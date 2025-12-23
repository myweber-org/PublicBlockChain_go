
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
	GetAllSessions() ([]Session, error)
	DeleteSession(sessionID string) error
}

type SessionCleaner struct {
	store SessionStore
}

func NewSessionCleaner(store SessionStore) *SessionCleaner {
	return &SessionCleaner{store: store}
}

func (sc *SessionCleaner) RunCleanup() error {
	sessions, err := sc.store.GetAllSessions()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, session := range sessions {
		if session.ExpiresAt.Before(now) {
			err := sc.store.DeleteSession(session.ID)
			if err != nil {
				log.Printf("Failed to delete session %s: %v", session.ID, err)
				continue
			}
			log.Printf("Deleted expired session: %s", session.ID)
		}
	}
	return nil
}

func (sc *SessionCleaner) StartDailyCleanup() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := sc.RunCleanup(); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			}
		}
	}
}

func main() {
	// In a real application, you would provide a proper SessionStore implementation
	var store SessionStore
	cleaner := NewSessionCleaner(store)
	
	// Run initial cleanup
	if err := cleaner.RunCleanup(); err != nil {
		log.Fatalf("Initial cleanup failed: %v", err)
	}
	
	// Start daily cleanup job
	cleaner.StartDailyCleanup()
}