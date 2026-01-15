package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidToken    = errors.New("invalid session token")
)

type Session struct {
	UserID    string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Manager struct {
	client     *redis.Client
	expiration time.Duration
}

func NewManager(redisAddr string, expiration time.Duration) (*Manager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Manager{
		client:     client,
		expiration: expiration,
	}, nil
}

func (m *Manager) Create(userID, username string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.expiration),
	}

	ctx := context.Background()
	err = m.client.Set(ctx, token, session, m.expiration).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *Manager) Validate(token string) (*Session, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	ctx := context.Background()
	var session Session
	err := m.client.Get(ctx, token).Scan(&session)
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		m.Delete(token)
		return nil, ErrSessionNotFound
	}

	return &session, nil
}

func (m *Manager) Delete(token string) error {
	ctx := context.Background()
	return m.client.Del(ctx, token).Err()
}

func (m *Manager) Refresh(token string) error {
	session, err := m.Validate(token)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(m.expiration)
	ctx := context.Background()
	return m.client.Set(ctx, token, session, m.expiration).Err()
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
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

    sessionID := generateSessionID()
    session := &Session{
        ID:        sessionID,
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(ttl),
    }
    sm.sessions[sessionID] = session
    return session
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

func (sm *SessionManager) DeleteSession(sessionID string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    delete(sm.sessions, sessionID)
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