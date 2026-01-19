
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
}package main

import (
	"context"
	"log"
	"time"
)

type SessionStore interface {
	DeleteExpiredSessions(ctx context.Context) error
}

type CleanupJob struct {
	store     SessionStore
	interval  time.Duration
}

func NewCleanupJob(store SessionStore, interval time.Duration) *CleanupJob {
	return &CleanupJob{
		store:    store,
		interval: interval,
	}
}

func (j *CleanupJob) Run(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Cleanup job stopped")
			return
		case <-ticker.C:
			if err := j.store.DeleteExpiredSessions(ctx); err != nil {
				log.Printf("Failed to delete expired sessions: %v", err)
			} else {
				log.Println("Expired sessions cleaned up successfully")
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := NewMemorySessionStore()
	job := NewCleanupJob(store, 24*time.Hour)

	go job.Run(ctx)

	<-ctx.Done()
}package main

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
	DeleteExpiredSessions(ctx context.Context, olderThan time.Time) error
}

type CleanupJob struct {
	store SessionStore
}

func NewCleanupJob(store SessionStore) *CleanupJob {
	return &CleanupJob{store: store}
}

func (j *CleanupJob) Run() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			j.cleanup()
		}
	}
}

func (j *CleanupJob) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cutoff := time.Now().Add(-sessionTTL)
	if err := j.store.DeleteExpiredSessions(ctx, cutoff); err != nil {
		log.Printf("Failed to clean up expired sessions: %v", err)
	} else {
		log.Printf("Successfully cleaned up sessions older than %v", cutoff)
	}
}

func main() {
	// In a real application, initialize your session store here
	// store := NewDatabaseSessionStore()
	// job := NewCleanupJob(store)
	// job.Run()
}