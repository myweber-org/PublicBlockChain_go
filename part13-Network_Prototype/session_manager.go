package session

import (
	"sync"
	"time"
)

type Session struct {
	ID        string
	Data      map[string]interface{}
	ExpiresAt time.Time
	mu        sync.RWMutex
}

type Manager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	stopChan chan struct{}
}

func NewManager(cleanupInterval time.Duration) *Manager {
	m := &Manager{
		sessions: make(map[string]*Session),
		stopChan: make(chan struct{}),
	}
	go m.startCleanup(cleanupInterval)
	return m
}

func (m *Manager) CreateSession(id string, ttl time.Duration) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := &Session{
		ID:        id,
		Data:      make(map[string]interface{}),
		ExpiresAt: time.Now().Add(ttl),
	}
	m.sessions[id] = session
	return session
}

func (m *Manager) GetSession(id string) *Session {
	m.mu.RLock()
	session, exists := m.sessions[id]
	m.mu.RUnlock()

	if !exists || time.Now().After(session.ExpiresAt) {
		return nil
	}
	return session
}

func (m *Manager) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanupExpired()
		case <-m.stopChan:
			return
		}
	}
}

func (m *Manager) cleanupExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, id)
		}
	}
}

func (m *Manager) Stop() {
	close(m.stopChan)
}

func (s *Session) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[key] = value
}

func (s *Session) Get(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Data[key]
}

func (s *Session) Extend(ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ExpiresAt = time.Now().Add(ttl)
}