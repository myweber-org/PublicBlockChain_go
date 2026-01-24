package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Endpoint  string    `json:"endpoint"`
    IPAddress string    `json:"ip_address"`
}

type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.RWMutex
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

func (rl *RateLimiter) Allow(key string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    windowStart := now.Add(-rl.window)

    timestamps := rl.requests[key]
    validRequests := []time.Time{}

    for _, ts := range timestamps {
        if ts.After(windowStart) {
            validRequests = append(validRequests, ts)
        }
    }

    if len(validRequests) >= rl.limit {
        return false
    }

    validRequests = append(validRequests, now)
    rl.requests[key] = validRequests
    return true
}

func loggingMiddleware(next http.Handler, logger *log.Logger, limiter *RateLimiter) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        clientIP := r.RemoteAddr
        if !limiter.Allow(clientIP) {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        start := time.Now()
        userID := r.Header.Get("X-User-ID")
        if userID == "" {
            userID = "anonymous"
        }

        activity := ActivityLog{
            Timestamp: time.Now(),
            UserID:    userID,
            Action:    r.Method,
            Endpoint:  r.URL.Path,
            IPAddress: clientIP,
        }

        logData, err := json.Marshal(activity)
        if err != nil {
            logger.Printf("Failed to marshal activity log: %v", err)
        } else {
            logger.Println(string(logData))
        }

        next.ServeHTTP(w, r)

        duration := time.Since(start)
        logger.Printf("Request completed in %v", duration)
    })
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{
        "status":  "success",
        "message": "Request processed successfully",
    }
    json.NewEncoder(w).Encode(response)
}

func main() {
    logger := log.New(log.Writer(), "ACTIVITY: ", log.LstdFlags)
    limiter := NewRateLimiter(100, time.Minute)

    mux := http.NewServeMux()
    mux.HandleFunc("/", mainHandler)

    wrappedHandler := loggingMiddleware(mux, logger, limiter)

    server := &http.Server{
        Addr:         ":8080",
        Handler:      wrappedHandler,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    logger.Println("Starting server on :8080")
    if err := server.ListenAndServe(); err != nil {
        logger.Fatal(err)
    }
}
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Path      string
	Method    string
	Timestamp time.Time
	IPAddress string
}

type ActivityLogger struct {
	activityChan chan ActivityLog
}

func NewActivityLogger(bufferSize int) *ActivityLogger {
	al := &ActivityLogger{
		activityChan: make(chan ActivityLog, bufferSize),
	}
	go al.processLogs()
	return al
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		activity := ActivityLog{
			UserID:    userID,
			Path:      r.URL.Path,
			Method:    r.Method,
			Timestamp: time.Now(),
			IPAddress: r.RemoteAddr,
		}

		select {
		case al.activityChan <- activity:
		default:
			log.Println("Activity log buffer full, dropping entry")
		}

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) processLogs() {
	for activity := range al.activityChan {
		log.Printf("Activity: User=%s %s %s from %s at %s",
			activity.UserID,
			activity.Method,
			activity.Path,
			activity.IPAddress,
			activity.Timestamp.Format(time.RFC3339))
	}
}

func (al *ActivityLogger) Close() {
	close(al.activityChan)
}