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
}