package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityEvent struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
    Details   string    `json:"details"`
}

type ActivityLogger struct {
    logFile *os.File
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &ActivityLogger{logFile: file}, nil
}

func (l *ActivityLogger) LogActivity(userID, eventType, details string) error {
    event := ActivityEvent{
        Timestamp: time.Now().UTC(),
        UserID:    userID,
        EventType: eventType,
        Details:   details,
    }

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    _, err = l.logFile.Write(append(data, '\n'))
    return err
}

func (l *ActivityLogger) Close() error {
    return l.logFile.Close()
}

func main() {
    logger, err := NewActivityLogger("activity.log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    err = logger.LogActivity("user123", "login", "User logged in from IP 192.168.1.100")
    if err != nil {
        log.Println("Failed to log activity:", err)
    }

    err = logger.LogActivity("user123", "file_upload", "Uploaded document.pdf")
    if err != nil {
        log.Println("Failed to log activity:", err)
    }

    fmt.Println("Activity logging completed")
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	IPAddress string
	Method    string
	Path      string
	Timestamp time.Time
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserID(authHeader)
		}

		activity := ActivityLog{
			UserID:    userID,
			IPAddress: r.RemoteAddr,
			Method:    r.Method,
			Path:      r.URL.Path,
			Timestamp: start,
		}

		log.Printf("Activity: %s %s by %s from %s", activity.Method, activity.Path, activity.UserID, activity.IPAddress)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(token string) string {
	return "user_" + token[:8]
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	al.handler.ServeHTTP(w, r)
	
	duration := time.Since(start)
	
	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	redisClient *redis.Client
	limiter     *rate.Limiter
	keyPrefix   string
}

func NewActivityLogger(redisAddr string, prefix string) *ActivityLogger {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	return &ActivityLogger{
		redisClient: rdb,
		limiter:     rate.NewLimiter(rate.Every(time.Minute), 10),
		keyPrefix:   prefix,
	}
}

func (al *ActivityLogger) LogActivity(ctx context.Context, userID string, action string) error {
	if !al.limiter.Allow() {
		return fmt.Errorf("rate limit exceeded for user %s", userID)
	}

	key := fmt.Sprintf("%s:activity:%s", al.keyPrefix, userID)
	timestamp := time.Now().Unix()

	activity := map[string]interface{}{
		"action":    action,
		"timestamp": timestamp,
		"user_id":   userID,
	}

	err := al.redisClient.HSet(ctx, key, activity).Err()
	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	expiration := 24 * time.Hour
	al.redisClient.Expire(ctx, key, expiration)

	return nil
}

func (al *ActivityLogger) GetRecentActivities(ctx context.Context, userID string) ([]map[string]string, error) {
	key := fmt.Sprintf("%s:activity:%s", al.keyPrefix, userID)
	result, err := al.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}

	var activities []map[string]string
	for field, value := range result {
		activity := map[string]string{
			"field": field,
			"value": value,
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (al *ActivityLogger) Close() error {
	return al.redisClient.Close()
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	userID := extractUserID(r)
	path := r.URL.Path
	method := r.Method

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("User %s %s %s completed in %v", userID, method, path, duration)
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return auth[:8]
	}
	return "anonymous"
}