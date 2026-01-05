
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter *RateLimiter
	logger      *log.Logger
}

type RateLimiter struct {
	requests map[string][]time.Time
	window   time.Duration
	maxReqs  int
}

func NewRateLimiter(window time.Duration, maxReqs int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		maxReqs:  maxReqs,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	timestamps := rl.requests[ip]

	var valid []time.Time
	for _, ts := range timestamps {
		if now.Sub(ts) <= rl.window {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= rl.maxReqs {
		return false
	}

	valid = append(valid, now)
	rl.requests[ip] = valid
	return true
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: NewRateLimiter(time.Minute, 100),
		logger:      logger,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if !al.rateLimiter.Allow(clientIP) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		defer func() {
			duration := time.Since(start)
			al.logger.Printf(
				"IP: %s | Method: %s | Path: %s | Status: %d | Duration: %v",
				clientIP,
				r.Method,
				r.URL.Path,
				recorder.statusCode,
				duration,
			)
		}()

		next.ServeHTTP(recorder, r)
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