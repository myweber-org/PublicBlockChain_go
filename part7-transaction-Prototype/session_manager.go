package session

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
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
    sessions map[string]*Session
    mu       sync.RWMutex
    duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
    return &Manager{
        sessions: make(map[string]*Session),
        duration: sessionDuration,
    }
}

func (m *Manager) Create(userID int) (string, error) {
    token, err := generateToken()
    if err != nil {
        return "", err
    }

    session := &Session{
        ID:        token,
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(m.duration),
    }

    m.mu.Lock()
    m.sessions[token] = session
    m.mu.Unlock()

    return token, nil
}

func (m *Manager) Validate(token string) (*Session, error) {
    m.mu.RLock()
    session, exists := m.sessions[token]
    m.mu.RUnlock()

    if !exists {
        return nil, errors.New("session not found")
    }

    if time.Now().After(session.ExpiresAt) {
        m.mu.Lock()
        delete(m.sessions, token)
        m.mu.Unlock()
        return nil, errors.New("session expired")
    }

    return session, nil
}

func (m *Manager) Refresh(token string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    session, exists := m.sessions[token]
    if !exists {
        return errors.New("session not found")
    }

    session.ExpiresAt = time.Now().Add(m.duration)
    return nil
}

func (m *Manager) Destroy(token string) {
    m.mu.Lock()
    delete(m.sessions, token)
    m.mu.Unlock()
}

func (m *Manager) Cleanup() {
    m.mu.Lock()
    defer m.mu.Unlock()

    now := time.Now()
    for token, session := range m.sessions {
        if now.After(session.ExpiresAt) {
            delete(m.sessions, token)
        }
    }
}

func generateToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}