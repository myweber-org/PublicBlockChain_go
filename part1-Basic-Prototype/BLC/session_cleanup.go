package main

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
	sessions []Session
}

func (s *SessionStore) RemoveExpiredSessions() {
	now := time.Now()
	var validSessions []Session

	for _, session := range s.sessions {
		if session.ExpiresAt.After(now) {
			validSessions = append(validSessions, session)
		}
	}

	s.sessions = validSessions
	log.Printf("Cleaned up expired sessions. Remaining: %d", len(s.sessions))
}

func main() {
	store := &SessionStore{
		sessions: []Session{
			{ID: "abc123", UserID: 1, ExpiresAt: time.Now().Add(-1 * time.Hour)},
			{ID: "def456", UserID: 2, ExpiresAt: time.Now().Add(24 * time.Hour)},
			{ID: "ghi789", UserID: 3, ExpiresAt: time.Now().Add(-5 * time.Minute)},
		},
	}

	store.RemoveExpiredSessions()
}