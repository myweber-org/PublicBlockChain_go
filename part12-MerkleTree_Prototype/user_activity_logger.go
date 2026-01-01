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
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		al.Logger.Printf(
			"[%s] %s %s %d %s %v",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			r.RemoteAddr,
			duration,
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
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	activities  map[string][]time.Time
	rateLimit   int
	window      time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		activities: make(map[string][]time.Time),
		rateLimit:  limit,
		window:     window,
	}
}

func (al *ActivityLogger) LogActivity(userID string, activity string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	key := userID + ":" + activity
	now := time.Now()

	// Clean old entries
	var validTimes []time.Time
	for _, t := range al.activities[key] {
		if now.Sub(t) <= al.window {
			validTimes = append(validTimes, t)
		}
	}

	if len(validTimes) >= al.rateLimit {
		return false
	}

	validTimes = append(validTimes, now)
	al.activities[key] = validTimes
	return true
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		activity := r.Method + ":" + r.URL.Path
		if !al.LogActivity(userID, activity) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}package middleware

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
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		al.Logger.Printf(
			"Method=%s Path=%s Status=%d Duration=%s RemoteAddr=%s UserAgent=%s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.RemoteAddr,
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
}