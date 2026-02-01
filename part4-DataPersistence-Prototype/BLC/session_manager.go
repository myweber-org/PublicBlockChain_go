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
}