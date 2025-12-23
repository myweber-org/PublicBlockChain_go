package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	Method    string
	Path      string
	UserAgent string
	IP        string
	Duration  time.Duration
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		activity := ActivityLog{
			Timestamp: start,
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IP:        getClientIP(r),
			Duration:  duration,
		}
		
		logActivity(activity, recorder.statusCode)
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

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func logActivity(activity ActivityLog, status int) {
	log.Printf("[ACTIVITY] %s | %s %s | IP: %s | Agent: %s | Duration: %v | Status: %d",
		activity.Timestamp.Format("2006-01-02 15:04:05"),
		activity.Method,
		activity.Path,
		activity.IP,
		activity.UserAgent,
		activity.Duration,
		status,
	)
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter *RateLimiter
}

type RateLimiter struct {
	requests map[string][]time.Time
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

func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	timestamps := rl.requests[ip]

	var validTimestamps []time.Time
	for _, ts := range timestamps {
		if now.Sub(ts) <= rl.window {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	if len(validTimestamps) >= rl.limit {
		return false
	}

	validTimestamps = append(validTimestamps, now)
	rl.requests[ip] = validTimestamps
	return true
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: NewRateLimiter(limit, window),
	}
}

func (al *ActivityLogger) LogActivity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		if !al.rateLimiter.Allow(clientIP) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			log.Printf("IP: %s | Method: %s | Path: %s | Duration: %v",
				clientIP, r.Method, r.URL.Path, duration)
		}()

		next.ServeHTTP(w, r)
	}
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request processed successfully"))
}
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
}