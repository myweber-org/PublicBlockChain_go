package main

import (
	"log"
	"time"
)

type SessionStore struct {
	sessions map[string]time.Time
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]time.Time),
	}
}

func (s *SessionStore) AddSession(id string) {
	s.sessions[id] = time.Now()
}

func (s *SessionStore) IsValidSession(id string) bool {
	created, exists := s.sessions[id]
	if !exists {
		return false
	}
	return time.Since(created) < 24*time.Hour
}

func (s *SessionStore) CleanupExpiredSessions() {
	cutoff := time.Now().Add(-24 * time.Hour)
	for id, created := range s.sessions {
		if created.Before(cutoff) {
			delete(s.sessions, id)
		}
	}
}

func startCleanupJob(store *SessionStore) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		store.CleanupExpiredSessions()
		log.Println("Expired sessions cleaned up")
	}
}

func main() {
	store := NewSessionStore()
	go startCleanupJob(store)

	store.AddSession("user123")
	log.Println("Session valid:", store.IsValidSession("user123"))

	time.Sleep(25 * time.Hour)
	log.Println("Session valid after 25h:", store.IsValidSession("user123"))
}