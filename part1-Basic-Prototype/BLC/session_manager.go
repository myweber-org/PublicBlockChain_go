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

func (sm *SessionManager) InvalidateSession(sessionID string) {
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
	return "session_" + time.Now().Format("20060102150405") + "_" + randomString(16)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}package session

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "time"

    "github.com/go-redis/redis/v8"
    "golang.org/x/net/context"
)

var (
    ErrInvalidToken = errors.New("invalid session token")
    ErrSessionExpired = errors.New("session has expired")
)

type Session struct {
    UserID    string
    Username  string
    CreatedAt time.Time
    ExpiresAt time.Time
}

type Manager struct {
    client    *redis.Client
    prefix    string
    expiry    time.Duration
}

func NewManager(client *redis.Client, prefix string, expiry time.Duration) *Manager {
    return &Manager{
        client: client,
        prefix: prefix,
        expiry: expiry,
    }
}

func (m *Manager) GenerateToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func (m *Manager) Create(ctx context.Context, userID, username string) (string, *Session, error) {
    token, err := m.GenerateToken()
    if err != nil {
        return "", nil, err
    }

    now := time.Now()
    session := &Session{
        UserID:    userID,
        Username:  username,
        CreatedAt: now,
        ExpiresAt: now.Add(m.expiry),
    }

    key := m.prefix + token
    err = m.client.Set(ctx, key, userID, m.expiry).Err()
    if err != nil {
        return "", nil, err
    }

    return token, session, nil
}

func (m *Manager) Validate(ctx context.Context, token string) (*Session, error) {
    if token == "" {
        return nil, ErrInvalidToken
    }

    key := m.prefix + token
    userID, err := m.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return nil, ErrInvalidToken
        }
        return nil, err
    }

    ttl, err := m.client.TTL(ctx, key).Result()
    if err != nil {
        return nil, err
    }

    now := time.Now()
    return &Session{
        UserID:    userID,
        CreatedAt: now.Add(-m.expiry + ttl),
        ExpiresAt: now.Add(ttl),
    }, nil
}

func (m *Manager) Refresh(ctx context.Context, token string) error {
    if token == "" {
        return ErrInvalidToken
    }

    key := m.prefix + token
    exists, err := m.client.Exists(ctx, key).Result()
    if err != nil {
        return err
    }
    if exists == 0 {
        return ErrInvalidToken
    }

    return m.client.Expire(ctx, key, m.expiry).Err()
}

func (m *Manager) Destroy(ctx context.Context, token string) error {
    if token == "" {
        return ErrInvalidToken
    }

    key := m.prefix + token
    return m.client.Del(ctx, key).Err()
}

func (m *Manager) CleanupExpired(ctx context.Context) error {
    keys, err := m.client.Keys(ctx, m.prefix+"*").Result()
    if err != nil {
        return err
    }

    for _, key := range keys {
        ttl, err := m.client.TTL(ctx, key).Result()
        if err != nil {
            continue
        }
        if ttl <= 0 {
            m.client.Del(ctx, key)
        }
    }
    return nil
}