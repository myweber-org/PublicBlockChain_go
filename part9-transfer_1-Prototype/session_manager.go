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
}package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type SessionManager struct {
	redisClient *redis.Client
	ttl         time.Duration
}

func NewSessionManager(redisAddr string, ttl time.Duration) (*SessionManager, error) {
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

	return &SessionManager{
		redisClient: client,
		ttl:         ttl,
	}, nil
}

func (sm *SessionManager) GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (sm *SessionManager) CreateSession(userID string) (string, error) {
	token, err := sm.GenerateToken()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	err = sm.redisClient.Set(ctx, token, userID, sm.ttl).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (sm *SessionManager) ValidateSession(token string) (string, error) {
	ctx := context.Background()
	userID, err := sm.redisClient.Get(ctx, token).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("session not found")
		}
		return "", err
	}

	err = sm.redisClient.Expire(ctx, token, sm.ttl).Err()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (sm *SessionManager) DestroySession(token string) error {
	ctx := context.Background()
	return sm.redisClient.Del(ctx, token).Err()
}

func (sm *SessionManager) CleanupExpiredSessions() error {
	return nil
}