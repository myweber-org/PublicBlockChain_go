
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
)

const cleanupInterval = 24 * time.Hour

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(ctx, db); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			} else {
				log.Println("Session cleanup completed successfully")
			}
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *database.DB) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("Cleaned up %d expired sessions", rowsAffected)
	return nil
}package main

import (
    "log"
    "time"
)

type Session struct {
    ID        string
    UserID    int
    ExpiresAt time.Time
}

type SessionStore struct {
    sessions map[string]Session
}

func NewSessionStore() *SessionStore {
    return &SessionStore{
        sessions: make(map[string]Session),
    }
}

func (s *SessionStore) CleanExpiredSessions() {
    now := time.Now()
    expiredCount := 0
    
    for id, session := range s.sessions {
        if session.ExpiresAt.Before(now) {
            delete(s.sessions, id)
            expiredCount++
        }
    }
    
    if expiredCount > 0 {
        log.Printf("Cleaned %d expired sessions", expiredCount)
    }
}

func startCleanupJob(store *SessionStore, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            store.CleanExpiredSessions()
        }
    }
}

func main() {
    store := NewSessionStore()
    
    // Add some test sessions
    store.sessions["abc123"] = Session{
        ID:        "abc123",
        UserID:    1,
        ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
    }
    
    store.sessions["def456"] = Session{
        ID:        "def456",
        UserID:    2,
        ExpiresAt: time.Now().Add(24 * time.Hour), // Valid for 24 hours
    }
    
    // Start cleanup job running every 5 minutes
    go startCleanupJob(store, 5*time.Minute)
    
    // Keep main goroutine alive
    select {}
}