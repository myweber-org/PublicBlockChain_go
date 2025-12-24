package main

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

func (sm *SessionManager) CreateSession(userID int, data map[string]interface{}) string {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    sessionID := generateSessionID()
    session := &Session{
        ID:        sessionID,
        UserID:    userID,
        Data:      data,
        ExpiresAt: time.Now().Add(sm.ttl),
    }
    sm.sessions[sessionID] = session
    return sessionID
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    session, exists := sm.sessions[sessionID]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil, false
    }
    return session, true
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
    stopChan chan struct{}
}

func NewSessionManager(cleanupInterval time.Duration) *SessionManager {
    sm := &SessionManager{
        sessions: make(map[string]*Session),
        stopChan: make(chan struct{}),
    }
    go sm.startCleanupRoutine(cleanupInterval)
    return sm
}

func (sm *SessionManager) CreateSession(userID int, ttl time.Duration) *Session {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    session := &Session{
        ID:        generateSessionID(),
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(ttl),
    }
    sm.sessions[session.ID] = session
    return session
}

func (sm *SessionManager) GetSession(id string) *Session {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    session, exists := sm.sessions[id]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil
    }
    return session
}

func (sm *SessionManager) InvalidateSession(id string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    delete(sm.sessions, id)
}

func (sm *SessionManager) startCleanupRoutine(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            sm.cleanupExpiredSessions()
        case <-sm.stopChan:
            return
        }
    }
}

func (sm *SessionManager) cleanupExpiredSessions() {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    now := time.Now()
    for id, session := range sm.sessions {
        if now.After(session.ExpiresAt) {
            delete(sm.sessions, id)
        }
    }
}

func (sm *SessionManager) Stop() {
    close(sm.stopChan)
}

func generateSessionID() string {
    return "sess_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
}