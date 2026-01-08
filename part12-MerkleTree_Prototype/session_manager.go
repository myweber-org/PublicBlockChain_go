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
}