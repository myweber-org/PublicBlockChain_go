
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
	ttl      time.Duration
}

func NewManager(ttl time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}
}

func (m *Manager) Create(userID int, data map[string]interface{}) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := &Session{
		ID:        token,
		UserID:    userID,
		Data:      data,
		ExpiresAt: time.Now().Add(m.ttl),
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
		m.Delete(token)
		return nil, errors.New("session expired")
	}

	return session, nil
}

func (m *Manager) Delete(token string) {
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
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	Data      map[string]interface{}
	ExpiresAt time.Time
}

type Manager struct {
	sessions map[string]Session
	mu       sync.RWMutex
	ttl      time.Duration
}

func NewManager(ttl time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]Session),
		ttl:      ttl,
	}
}

func (m *Manager) Create(userID string, data map[string]interface{}) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		ID:        token,
		UserID:    userID,
		Data:      data,
		ExpiresAt: time.Now().Add(m.ttl),
	}

	m.mu.Lock()
	m.sessions[token] = session
	m.mu.Unlock()

	return token, nil
}

func (m *Manager) Validate(token string) (Session, error) {
	m.mu.RLock()
	session, exists := m.sessions[token]
	m.mu.RUnlock()

	if !exists {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		m.Delete(token)
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func (m *Manager) Delete(token string) {
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
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session struct {
	UserID    string
	ExpiresAt time.Time
	Data      map[string]interface{}
}

type Manager struct {
	sessions map[string]Session
	mu       sync.RWMutex
	ttl      time.Duration
}

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

func NewManager(ttl time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]Session),
		ttl:      ttl,
	}
}

func (m *Manager) Create(userID string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(m.ttl),
		Data:      make(map[string]interface{}),
	}

	m.mu.Lock()
	m.sessions[token] = session
	m.mu.Unlock()

	return token, nil
}

func (m *Manager) Validate(token string) (Session, error) {
	m.mu.RLock()
	session, exists := m.sessions[token]
	m.mu.RUnlock()

	if !exists {
		return Session{}, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		m.mu.Lock()
		delete(m.sessions, token)
		m.mu.Unlock()
		return Session{}, ErrSessionExpired
	}

	return session, nil
}

func (m *Manager) Refresh(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[token]
	if !exists {
		return ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		delete(m.sessions, token)
		return ErrSessionExpired
	}

	session.ExpiresAt = time.Now().Add(m.ttl)
	m.sessions[token] = session
	return nil
}

func (m *Manager) Delete(token string) {
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

type Manager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	ttl      time.Duration
}

func NewManager(ttl time.Duration) *Manager {
	m := &Manager{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}
	go m.cleanupLoop()
	return m
}

func (m *Manager) Create(userID int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := &Session{
		ID:        generateID(),
		UserID:    userID,
		Data:      make(map[string]interface{}),
		ExpiresAt: time.Now().Add(m.ttl),
	}
	m.sessions[session.ID] = session
	return session
}

func (m *Manager) Get(id string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[id]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil
	}
	return session
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
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

func generateID() string {
	return "sess_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}