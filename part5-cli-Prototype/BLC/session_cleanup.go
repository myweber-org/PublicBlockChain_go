package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	sessionKeyPattern = "session:*"
	batchSize         = 100
)

func cleanupExpiredSessions(rdb *redis.Client) error {
	ctx := context.Background()
	var cursor uint64
	var keys []string
	var err error

	for {
		keys, cursor, err = rdb.Scan(ctx, cursor, sessionKeyPattern, batchSize).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			ttl, err := rdb.TTL(ctx, key).Result()
			if err != nil {
				log.Printf("Failed to get TTL for key %s: %v", key, err)
				continue
			}
			if ttl < 0 {
				if err := rdb.Del(ctx, key).Err(); err != nil {
					log.Printf("Failed to delete expired key %s: %v", key, err)
				} else {
					log.Printf("Deleted expired session: %s", key)
				}
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Starting session cleanup...")
		if err := cleanupExpiredSessions(rdb); err != nil {
			log.Printf("Session cleanup failed: %v", err)
		} else {
			log.Println("Session cleanup completed")
		}
	}
}
package main

import (
	"log"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type SessionStore interface {
	GetAllSessions() ([]Session, error)
	DeleteSession(id string) error
}

type DBSessionStore struct{}

func (s *DBSessionStore) GetAllSessions() ([]Session, error) {
	return []Session{}, nil
}

func (s *DBSessionStore) DeleteSession(id string) error {
	return nil
}

func cleanupExpiredSessions(store SessionStore) error {
	sessions, err := store.GetAllSessions()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, session := range sessions {
		if session.ExpiresAt.Before(now) {
			if err := store.DeleteSession(session.ID); err != nil {
				log.Printf("Failed to delete session %s: %v", session.ID, err)
			} else {
				log.Printf("Deleted expired session: %s", session.ID)
			}
		}
	}
	return nil
}

func scheduleCleanup(store SessionStore, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := cleanupExpiredSessions(store); err != nil {
			log.Printf("Session cleanup failed: %v", err)
		}
	}
}

func main() {
	store := &DBSessionStore{}
	go scheduleCleanup(store, 24*time.Hour)

	select {}
}