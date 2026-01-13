
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
    mu        sync.RWMutex
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
    sessionID := generateSessionID()
    session := &Session{
        ID:        sessionID,
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(ttl),
    }

    sm.mu.Lock()
    sm.sessions[sessionID] = session
    sm.mu.Unlock()

    return session
}

func (sm *SessionManager) GetSession(sessionID string) *Session {
    sm.mu.RLock()
    session, exists := sm.sessions[sessionID]
    sm.mu.RUnlock()

    if !exists || time.Now().After(session.ExpiresAt) {
        return nil
    }
    return session
}

func (sm *SessionManager) InvalidateSession(sessionID string) {
    sm.mu.Lock()
    delete(sm.sessions, sessionID)
    sm.mu.Unlock()
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

func generateSessionID() string {
    return "session_" + time.Now().Format("20060102150405") + "_" + randomString(16)
}

func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
}
package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionManager struct {
	sessions map[string]Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]Session),
	}
}

func (sm *SessionManager) CreateSession(userID int) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		ID:        token,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	sm.sessions[token] = session
	return token, nil
}

func (sm *SessionManager) ValidateSession(token string) (Session, error) {
	session, exists := sm.sessions[token]
	if !exists {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(sm.sessions, token)
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func (sm *SessionManager) InvalidateSession(token string) {
	delete(sm.sessions, token)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}