package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	userLimits  map[string]time.Time
	rateLimit   time.Duration
	logFilePath string
}

func NewActivityLogger(limit time.Duration, logFile string) *ActivityLogger {
	return &ActivityLogger{
		userLimits:  make(map[string]time.Time),
		rateLimit:   limit,
		logFilePath: logFile,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		now := time.Now()
		al.mu.RLock()
		lastLog, exists := al.userLimits[userID]
		al.mu.RUnlock()

		if exists && now.Sub(lastLog) < al.rateLimit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		al.mu.Lock()
		al.userLimits[userID] = now
		al.mu.Unlock()

		logEntry := struct {
			Timestamp time.Time `json:"timestamp"`
			UserID    string    `json:"user_id"`
			Method    string    `json:"method"`
			Path      string    `json:"path"`
			IP        string    `json:"ip"`
		}{
			Timestamp: now,
			UserID:    userID,
			Method:    r.Method,
			Path:      r.URL.Path,
			IP:        r.RemoteAddr,
		}

		log.Printf("Activity: %s %s by %s from %s", logEntry.Method, logEntry.Path, logEntry.UserID, logEntry.IP)

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) CleanupOldEntries() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			al.mu.Lock()
			cutoff := time.Now().Add(-24 * time.Hour)
			for userID, lastLog := range al.userLimits {
				if lastLog.Before(cutoff) {
					delete(al.userLimits, userID)
				}
			}
			al.mu.Unlock()
		}
	}()
}