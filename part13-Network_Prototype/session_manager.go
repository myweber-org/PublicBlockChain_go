
package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	Data      map[string]interface{}
}

type Manager struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	ttl      time.Duration
}

func NewManager(ttl time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}
}

func (m *Manager) CreateSession(userID string, data map[string]interface{}) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := &Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(m.ttl),
		Data:      data,
	}

	m.mutex.Lock()
	m.sessions[token] = session
	m.mutex.Unlock()

	return token, nil
}

func (m *Manager) ValidateSession(token string) (*Session, error) {
	m.mutex.RLock()
	session, exists := m.sessions[token]
	m.mutex.RUnlock()

	if !exists {
		return nil, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		m.DeleteSession(token)
		return nil, errors.New("session expired")
	}

	return session, nil
}

func (m *Manager) DeleteSession(token string) {
	m.mutex.Lock()
	delete(m.sessions, token)
	m.mutex.Unlock()
}

func (m *Manager) CleanupExpired() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for token, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, token)
		}
	}
}

func (m *Manager) RefreshSession(token string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	session, exists := m.sessions[token]
	if !exists {
		return errors.New("session not found")
	}

	session.ExpiresAt = time.Now().Add(m.ttl)
	return nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}