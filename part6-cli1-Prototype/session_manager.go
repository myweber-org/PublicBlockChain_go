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

type Manager struct {
	sessions map[string]Session
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]Session),
	}
}

func (m *Manager) CreateSession(userID int, duration time.Duration) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		ID:        token,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	m.sessions[token] = session
	return token, nil
}

func (m *Manager) ValidateSession(token string) (Session, error) {
	session, exists := m.sessions[token]
	if !exists {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(m.sessions, token)
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func (m *Manager) InvalidateSession(token string) {
	delete(m.sessions, token)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}