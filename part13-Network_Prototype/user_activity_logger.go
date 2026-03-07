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
				log.Printf("Rate limited: %s", userIP)
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
		}

		al.rateLimiter[userIP] = now

		log.Printf("Activity: %s %s from %s", r.Method, r.URL.Path, userIP)

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) Cleanup() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			now := time.Now()
			for ip, lastSeen := range al.rateLimiter {
				if now.Sub(lastSeen) > 24*time.Hour {
					delete(al.rateLimiter, ip)
				}
			}
		}
	}()
}