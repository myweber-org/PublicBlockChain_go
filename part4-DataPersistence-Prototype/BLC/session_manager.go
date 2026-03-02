package session

import (
    "sync"
    "time"
)

type Session struct {
    ID        string
    UserID    int
    Data      map[string]interface{}
    ExpiresAt time.Time
}

type Manager struct {
    sessions map[string]Session
    mu       sync.RWMutex
    duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
    m := &Manager{
        sessions: make(map[string]Session),
        duration: sessionDuration,
    }
    go m.cleanupRoutine()
    return m
}

func (m *Manager) Create(userID int) Session {
    m.mu.Lock()
    defer m.mu.Unlock()

    session := Session{
        ID:        generateUUID(),
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(m.duration),
    }
    m.sessions[session.ID] = session
    return session
}

func (m *Manager) Get(sessionID string) (Session, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    session, exists := m.sessions[sessionID]
    if !exists || time.Now().After(session.ExpiresAt) {
        return Session{}, false
    }
    return session, true
}

func (m *Manager) Refresh(sessionID string) bool {
    m.mu.Lock()
    defer m.mu.Unlock()

    session, exists := m.sessions[sessionID]
    if !exists {
        return false
    }
    session.ExpiresAt = time.Now().Add(m.duration)
    m.sessions[sessionID] = session
    return true
}

func (m *Manager) cleanupRoutine() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        m.mu.Lock()
        now := time.Now()
        for id, session := range m.sessions {
            if now.After(session.ExpiresAt) {
                delete(m.sessions, id)
            }
        }
        m.mu.Unlock()
    }
}

func generateUUID() string {
    return "generated-unique-id"
}package main

import (
    "sync"
    "time"
)

type Session struct {
    ID        string
    UserID    int
    Data      map[string]interface{}
    ExpiresAt time.Time
}

type SessionManager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
    ttl      time.Duration
}

func NewSessionManager(ttl time.Duration) *SessionManager {
    sm := &SessionManager{
        sessions: make(map[string]*Session),
        ttl:      ttl,
    }
    go sm.cleanupWorker()
    return sm
}

func (sm *SessionManager) CreateSession(userID int) *Session {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    sessionID := generateSessionID()
    session := &Session{
        ID:        sessionID,
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(sm.ttl),
    }
    sm.sessions[sessionID] = session
    return session
}

func (sm *SessionManager) GetSession(sessionID string) *Session {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    session, exists := sm.sessions[sessionID]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil
    }
    return session
}

func (sm *SessionManager) RefreshSession(sessionID string) bool {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    session, exists := sm.sessions[sessionID]
    if !exists || time.Now().After(session.ExpiresAt) {
        return false
    }
    session.ExpiresAt = time.Now().Add(sm.ttl)
    return true
}

func (sm *SessionManager) cleanupWorker() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        sm.mu.Lock()
        now := time.Now()
        for id, session := range sm.sessions {
            if now.After(session.ExpiresAt) {
                delete(sm.sessions, id)
            }
        }
        sm.mu.Unlock()
    }
}

func generateSessionID() string {
    return "sess_" + time.Now().Format("20060102150405") + "_" + randomString(16)
}

func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
}package session

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

	if !exists {
		return nil
	}

	session.mu.RLock()
	defer session.mu.RUnlock()

	if time.Now().After(session.ExpiresAt) {
		go m.DeleteSession(id)
		return nil
	}
	return session
}

func (m *Manager) DeleteSession(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, id)
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
		session.mu.RLock()
		expired := now.After(session.ExpiresAt)
		session.mu.RUnlock()

		if expired {
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