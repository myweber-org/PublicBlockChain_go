package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	Logger *log.Logger
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{Logger: logger}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(startTime)
		
		al.Logger.Printf(
			"Method: %s | Path: %s | Status: %d | Duration: %v | User-Agent: %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.UserAgent(),
		)
	})
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

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"%s %s %s %d %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			rw.statusCode,
			duration,
		)
	})
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
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Details   string    `json:"details"`
}

type ActivityLogger struct {
    logFile string
}

func NewActivityLogger(logFile string) *ActivityLogger {
    return &ActivityLogger{logFile: logFile}
}

func (l *ActivityLogger) LogActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    file, err := os.OpenFile(l.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(logEntry); err != nil {
        return fmt.Errorf("failed to encode log entry: %w", err)
    }

    return nil
}

func main() {
    logger := NewActivityLogger("user_activity.log")

    activities := []struct {
        userID  string
        action  string
        details string
    }{
        {"user123", "LOGIN", "User logged in from IP 192.168.1.100"},
        {"user456", "UPLOAD", "File 'report.pdf' uploaded successfully"},
        {"user123", "SEARCH", "Searched for 'quarterly results'"},
        {"user789", "LOGOUT", "User session terminated"},
    }

    for _, activity := range activities {
        if err := logger.LogActivity(activity.userID, activity.action, activity.details); err != nil {
            fmt.Printf("Error logging activity: %v\n", err)
        } else {
            fmt.Printf("Logged: %s - %s\n", activity.userID, activity.action)
        }
    }
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

func NewActivityLogger(redisAddr string, prefix string) (*ActivityLogger, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &ActivityLogger{
		redisClient: rdb,
		limiter:     rate.NewLimiter(rate.Every(time.Minute), 100),
		keyPrefix:   prefix,
	}, nil
}

func (al *ActivityLogger) LogActivity(userID string, action string, metadata map[string]interface{}) error {
	if !al.limiter.Allow() {
		return fmt.Errorf("rate limit exceeded for user %s", userID)
	}

	ctx := context.Background()
	key := fmt.Sprintf("%s:activity:%s:%d", al.keyPrefix, userID, time.Now().UnixNano())

	data := map[string]interface{}{
		"user_id":  userID,
		"action":   action,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"metadata": metadata,
	}

	if err := al.redisClient.HSet(ctx, key, data).Err(); err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	expireKey := fmt.Sprintf("%s:activity:expiry:%s", al.keyPrefix, userID)
	al.redisClient.SAdd(ctx, expireKey, key)
	al.redisClient.Expire(ctx, expireKey, 30*24*time.Hour)

	return nil
}

func (al *ActivityLogger) GetRecentActivities(userID string, limit int64) ([]map[string]string, error) {
	ctx := context.Background()
	pattern := fmt.Sprintf("%s:activity:%s:*", al.keyPrefix, userID)

	keys, err := al.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch activity keys: %w", err)
	}

	if int64(len(keys)) > limit {
		keys = keys[:limit]
	}

	var activities []map[string]string
	for _, key := range keys {
		result, err := al.redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			continue
		}
		activities = append(activities, result)
	}

	return activities, nil
}

func (al *ActivityLogger) Close() error {
	return al.redisClient.Close()
}