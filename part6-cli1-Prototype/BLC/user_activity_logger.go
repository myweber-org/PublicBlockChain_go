package middleware

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
		"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s",
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
}
package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Endpoint  string    `json:"endpoint"`
	Method    string    `json:"method"`
	IPAddress string    `json:"ip_address"`
}

type RateLimiter struct {
	mu       sync.Mutex
	counters map[string]int
	resetAt  time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		counters: make(map[string]int),
		resetAt:  time.Now().Add(time.Hour),
	}
}

func (rl *RateLimiter) Allow(userID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if time.Now().After(rl.resetAt) {
		rl.counters = make(map[string]int)
		rl.resetAt = time.Now().Add(time.Hour)
	}

	if rl.counters[userID] >= 100 {
		return false
	}

	rl.counters[userID]++
	return true
}

func ActivityLogger(rateLimiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				userID = "anonymous"
			}

			if !rateLimiter.Allow(userID) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			activity := ActivityLog{
				Timestamp: time.Now().UTC(),
				UserID:    userID,
				Endpoint:  r.URL.Path,
				Method:    r.Method,
				IPAddress: r.RemoteAddr,
			}

			logData, err := json.Marshal(activity)
			if err == nil {
				go func() {
					println(string(logData))
				}()
			}

			next.ServeHTTP(w, r)
		})
	}
}