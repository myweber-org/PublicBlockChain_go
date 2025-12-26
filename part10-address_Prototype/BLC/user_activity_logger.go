package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter *RateLimiter
}

func NewActivityLogger(requestsPerMinute int) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: NewRateLimiter(requestsPerMinute, time.Minute),
	}
}

func (al *ActivityLogger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !al.rateLimiter.Allow() {
			log.Printf("Rate limit exceeded for logging, skipping activity log")
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		userID := extractUserID(r)
		path := r.URL.Path
		method := r.Method

		defer func() {
			duration := time.Since(start)
			log.Printf("User %s %s %s completed in %v", userID, method, path, duration)
		}()

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if user := r.Header.Get("X-User-ID"); user != "" {
		return user
	}
	return "anonymous"
}

type RateLimiter struct {
	requests int
	interval time.Duration
	bucket   chan struct{}
}

func NewRateLimiter(requests int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: requests,
		interval: interval,
		bucket:   make(chan struct{}, requests),
	}

	for i := 0; i < requests; i++ {
		rl.bucket <- struct{}{}
	}

	go rl.refill()

	return rl
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.bucket:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) refill() {
	ticker := time.NewTicker(rl.interval / time.Duration(rl.requests))
	defer ticker.Stop()

	for range ticker.C {
		select {
		case rl.bucket <- struct{}{}:
		default:
		}
	}
}