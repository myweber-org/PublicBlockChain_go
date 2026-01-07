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
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityRecord struct {
	UserID    string
	IPAddress string
	Endpoint  string
	Method    string
	Timestamp time.Time
}

type ActivityLogger struct {
	mu      sync.RWMutex
	records map[string][]ActivityRecord
	limiter *RateLimiter
}

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if _, exists := rl.requests[key]; !exists {
		rl.requests[key] = []time.Time{}
	}

	validWindow := now.Add(-rl.window)
	validRequests := []time.Time{}
	for _, t := range rl.requests[key] {
		if t.After(validWindow) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= rl.limit {
		return false
	}

	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests
	return true
}

func NewActivityLogger() *ActivityLogger {
	return &ActivityLogger{
		records: make(map[string][]ActivityRecord),
		limiter: NewRateLimiter(100, time.Minute*5),
	}
}

func (al *ActivityLogger) LogActivity(userID, ip, endpoint, method string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	record := ActivityRecord{
		UserID:    userID,
		IPAddress: ip,
		Endpoint:  endpoint,
		Method:    method,
		Timestamp: time.Now(),
	}

	if _, exists := al.records[userID]; !exists {
		al.records[userID] = []ActivityRecord{}
	}
	al.records[userID] = append(al.records[userID], record)

	if len(al.records[userID]) > 1000 {
		al.records[userID] = al.records[userID][len(al.records[userID])-1000:]
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := "anonymous"
		if id, ok := ctx.Value("userID").(string); ok {
			userID = id
		}

		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		if !al.limiter.Allow(userID + ":" + ip) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		al.LogActivity(userID, ip, r.URL.Path, r.Method)

		log.Printf("Activity: %s %s %s %s", userID, ip, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) GetUserActivities(userID string) []ActivityRecord {
	al.mu.RLock()
	defer al.mu.RUnlock()

	if activities, exists := al.records[userID]; exists {
		return append([]ActivityRecord{}, activities...)
	}
	return []ActivityRecord{}
}

func (al *ActivityLogger) CleanupOldRecords(maxAge time.Duration) {
	al.mu.Lock()
	defer al.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	for userID, records := range al.records {
		filtered := []ActivityRecord{}
		for _, record := range records {
			if record.Timestamp.After(cutoff) {
				filtered = append(filtered, record)
			}
		}
		if len(filtered) == 0 {
			delete(al.records, userID)
		} else {
			al.records[userID] = filtered
		}
	}
}