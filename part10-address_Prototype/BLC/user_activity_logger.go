package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter map[string]time.Time
	window      time.Duration
}

func NewActivityLogger(window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: make(map[string]time.Time),
		window:      window,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIP := r.RemoteAddr
		now := time.Now()

		if lastSeen, exists := al.rateLimiter[userIP]; exists {
			if now.Sub(lastSeen) < al.window {
				log.Printf("Rate limit exceeded for IP: %s", userIP)
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
		}

		al.rateLimiter[userIP] = now

		log.Printf("Activity from %s: %s %s", userIP, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) CleanupOldEntries() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for ip, lastSeen := range al.rateLimiter {
			if now.Sub(lastSeen) > 24*time.Hour {
				delete(al.rateLimiter, ip)
			}
		}
	}
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

	log.Printf("Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}