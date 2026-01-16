package session

import (
    "sync"
    "time"
)

type Session struct {
    ID        string
    Data      map[string]interface{}
    ExpiresAt time.Time
}

type Manager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
    duration time.Duration
}

func NewManager(duration time.Duration) *Manager {
    m := &Manager{
        sessions: make(map[string]*Session),
        duration: duration,
    }
    go m.cleanupRoutine()
    return m
}

func (m *Manager) Create(id string) *Session {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    session := &Session{
        ID:        id,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(m.duration),
    }
    m.sessions[id] = session
    return session
}

func (m *Manager) Get(id string) (*Session, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    session, exists := m.sessions[id]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil, false
    }
    return session, true
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
}package session

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

var sessions = make(map[string]Session)

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(userID int, duration time.Duration) (Session, error) {
	token, err := GenerateToken()
	if err != nil {
		return Session{}, err
	}

	session := Session{
		ID:        token,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	sessions[token] = session
	return session, nil
}

func ValidateSession(token string) (Session, error) {
	session, exists := sessions[token]
	if !exists {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(sessions, token)
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func DeleteSession(token string) {
	delete(sessions, token)
}package session

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

func (sm *SessionManager) Create(userID int) *Session {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    session := &Session{
        ID:        generateSessionID(),
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(sm.ttl),
    }
    sm.sessions[session.ID] = session
    return session
}

func (sm *SessionManager) Get(sessionID string) (*Session, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    session, exists := sm.sessions[sessionID]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil, false
    }
    return session, true
}

func (sm *SessionManager) Delete(sessionID string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    delete(sm.sessions, sessionID)
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