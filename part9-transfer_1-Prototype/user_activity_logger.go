package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	rateLimiter map[string][]time.Time
	windowSize  time.Duration
	maxRequests int
}

func NewActivityLogger(window time.Duration, max int) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: make(map[string][]time.Time),
		windowSize:  window,
		maxRequests: max,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		path := r.URL.Path

		if !al.checkRateLimit(clientIP) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("Activity: IP=%s, Agent=%s, Path=%s, Duration=%v", 
			clientIP, userAgent, path, duration)
	})
}

func (al *ActivityLogger) checkRateLimit(clientIP string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.windowSize)

	if _, exists := al.rateLimiter[clientIP]; !exists {
		al.rateLimiter[clientIP] = []time.Time{}
	}

	requests := al.rateLimiter[clientIP]
	var validRequests []time.Time
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= al.maxRequests {
		return false
	}

	validRequests = append(validRequests, now)
	al.rateLimiter[clientIP] = validRequests
	return true
}

func (al *ActivityLogger) CleanupOldEntries() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			al.mu.Lock()
			windowStart := time.Now().Add(-al.windowSize)
			for ip, requests := range al.rateLimiter {
				var validRequests []time.Time
				for _, t := range requests {
					if t.After(windowStart) {
						validRequests = append(validRequests, t)
					}
				}
				if len(validRequests) == 0 {
					delete(al.rateLimiter, ip)
				} else {
					al.rateLimiter[ip] = validRequests
				}
			}
			al.mu.Unlock()
		}
	}()
}