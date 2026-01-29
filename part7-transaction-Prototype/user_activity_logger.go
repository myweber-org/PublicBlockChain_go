package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu      sync.RWMutex
	entries map[string][]time.Time
	limit   int
	window  time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		entries: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (al *ActivityLogger) LogActivity(userID string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	// Clean old entries
	var validEntries []time.Time
	for _, t := range al.entries[userID] {
		if t.After(windowStart) {
			validEntries = append(validEntries, t)
		}
	}

	if len(validEntries) >= al.limit {
		return false
	}

	validEntries = append(validEntries, now)
	al.entries[userID] = validEntries
	return true
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !al.LogActivity(userID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		log.Printf("Activity logged for user %s: %s %s", userID, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) Cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		al.mu.Lock()
		windowStart := time.Now().Add(-al.window)
		for userID, times := range al.entries {
			var validEntries []time.Time
			for _, t := range times {
				if t.After(windowStart) {
					validEntries = append(validEntries, t)
				}
			}
			if len(validEntries) == 0 {
				delete(al.entries, userID)
			} else {
				al.entries[userID] = validEntries
			}
		}
		al.mu.Unlock()
	}
}