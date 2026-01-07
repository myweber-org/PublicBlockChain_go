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
}