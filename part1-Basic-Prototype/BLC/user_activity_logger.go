package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details,omitempty"`
}

type ActivityLogger struct {
	logFile *os.File
	encoder *json.Encoder
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ActivityLogger{
		logFile: file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (l *ActivityLogger) LogActivity(userID, action, details string) error {
	entry := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}
	return l.encoder.Encode(entry)
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

	activities := []struct {
		userID  string
		action  string
		details string
	}{
		{"user_001", "LOGIN", "Successful authentication"},
		{"user_002", "UPLOAD", "File: report.pdf"},
		{"user_001", "LOGOUT", "Session ended"},
	}

	for _, act := range activities {
		if err := logger.LogActivity(act.userID, act.action, act.details); err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}

	fmt.Println("Activity logging completed")
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
	ipAddress := r.RemoteAddr

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("User %s from %s accessed %s %s - Duration: %v",
		userID, ipAddress, r.Method, r.URL.Path, duration)
}

func extractUserID(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "anonymous"
	}
	return authHeader[:min(8, len(authHeader))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(start)
	log.Printf(
		"Method: %s | Path: %s | Status: %d | Duration: %v | UserAgent: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.UserAgent(),
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
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
	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.RemoteAddr

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("Activity: %s %s | IP: %s | Agent: %s | Duration: %v",
		r.Method, r.URL.Path, ipAddress, userAgent, duration)
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter *RateLimiter
}

type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: &RateLimiter{
			requests: make(map[string][]time.Time),
			limit:    limit,
			window:   window,
		},
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		path := r.URL.Path

		if !al.rateLimiter.Allow(clientIP) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("Activity: IP=%s, UA=%s, Path=%s, Duration=%v", clientIP, userAgent, path, duration)
	})
}

func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	windowStart := now.Add(-rl.window)

	if requests, exists := rl.requests[ip]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[ip] = validRequests

		if len(validRequests) >= rl.limit {
			return false
		}
	}

	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}package middleware

import (
    "context"
    "net/http"
    "time"
)

type ActivityKey string

const (
    UserIDKey ActivityKey = "userID"
    ActionKey ActivityKey = "action"
)

type ActivityLogger interface {
    LogActivity(ctx context.Context, userID string, action string, timestamp time.Time)
}

type activityLogger struct {
    store ActivityStore
}

func NewActivityLogger(store ActivityStore) ActivityLogger {
    return &activityLogger{store: store}
}

func (al *activityLogger) LogActivity(ctx context.Context, userID string, action string, timestamp time.Time) {
    al.store.Save(ctx, userID, action, timestamp)
}

func ActivityLoggingMiddleware(logger ActivityLogger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        userID := r.Header.Get("X-User-ID")
        if userID != "" {
            ctx = context.WithValue(ctx, UserIDKey, userID)
        }
        
        action := r.Method + " " + r.URL.Path
        ctx = context.WithValue(ctx, ActionKey, action)
        
        start := time.Now()
        next.ServeHTTP(w, r.WithContext(ctx))
        
        if userID != "" {
            logger.LogActivity(ctx, userID, action, start)
        }
    })
}

type ActivityStore interface {
    Save(ctx context.Context, userID string, action string, timestamp time.Time)
}

type MemoryActivityStore struct {
    activities []ActivityRecord
}

type ActivityRecord struct {
    UserID    string
    Action    string
    Timestamp time.Time
}

func NewMemoryActivityStore() *MemoryActivityStore {
    return &MemoryActivityStore{
        activities: make([]ActivityRecord, 0),
    }
}

func (mas *MemoryActivityStore) Save(ctx context.Context, userID string, action string, timestamp time.Time) {
    record := ActivityRecord{
        UserID:    userID,
        Action:    action,
        Timestamp: timestamp,
    }
    mas.activities = append(mas.activities, record)
}

func (mas *MemoryActivityStore) GetActivities() []ActivityRecord {
    return mas.activities
}