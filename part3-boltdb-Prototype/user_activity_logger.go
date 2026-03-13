package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu      sync.RWMutex
	clients map[string][]time.Time
	limit   int
	window  time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		now := time.Now()

		al.mu.Lock()
		defer al.mu.Unlock()

		al.cleanupOldEntries(clientIP, now)

		if len(al.clients[clientIP]) >= al.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			log.Printf("Rate limit exceeded for client: %s", clientIP)
			return
		}

		al.clients[clientIP] = append(al.clients[clientIP], now)
		log.Printf("Activity logged: %s - %s %s", clientIP, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) cleanupOldEntries(clientIP string, now time.Time) {
	if entries, exists := al.clients[clientIP]; exists {
		validEntries := []time.Time{}
		for _, t := range entries {
			if now.Sub(t) <= al.window {
				validEntries = append(validEntries, t)
			}
		}
		if len(validEntries) == 0 {
			delete(al.clients, clientIP)
		} else {
			al.clients[clientIP] = validEntries
		}
	}
}

func (al *ActivityLogger) GetActivityCount(clientIP string) int {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return len(al.clients[clientIP])
}