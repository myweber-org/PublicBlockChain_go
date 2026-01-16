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

func activityLogger(next http.Handler) http.Handler {
    rl := NewRateLimiter(100, time.Minute)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        clientIP := r.RemoteAddr
        if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
            clientIP = forwarded
        }

        userID := r.Header.Get("X-User-ID")
        if userID == "" {
            userID = "anonymous"
        }

        logEntry := ActivityLog{
            Timestamp: time.Now().UTC(),
            UserID:    userID,
            Action:    r.Method,
            Endpoint:  r.URL.Path,
            IPAddress: clientIP,
        }

        if !rl.Allow(userID) {
            http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
            logData, _ := json.Marshal(logEntry)
            log.Printf("RATE_LIMIT: %s", string(logData))
            return
        }

        logData, err := json.Marshal(logEntry)
        if err != nil {
            log.Printf("Failed to marshal log entry: %v", err)
        } else {
            log.Printf("ACTIVITY: %s", string(logData))
        }

        next.ServeHTTP(w, r)
    })
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{
        "status":  "success",
        "message": "request processed",
    }
    json.NewEncoder(w).Encode(response)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", mainHandler)

    wrappedMux := activityLogger(mux)

    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
        log.Fatal(err)
    }
}package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ActivityLog struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address"`
}

type RateLimiter struct {
	mu       sync.Mutex
	counters map[string]int
	window   time.Duration
	limit    int
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		counters: make(map[string]int),
		window:   window,
		limit:    limit,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	count, exists := rl.counters[key]
	currentTime := time.Now()

	if !exists || time.Since(time.Unix(int64(count>>32), 0)) > rl.window {
		rl.counters[key] = (int(currentTime.Unix()) << 32) | 1
		return true
	}

	requests := count & 0xFFFFFFFF
	if requests >= rl.limit {
		return false
	}

	rl.counters[key] = (count & 0xFFFFFFFF00000000) | (requests + 1)
	return true
}

type ActivityLogger struct {
	rateLimiter *RateLimiter
	logChannel  chan ActivityLog
}

func NewActivityLogger() *ActivityLogger {
	logger := &ActivityLogger{
		rateLimiter: NewRateLimiter(time.Minute, 100),
		logChannel:  make(chan ActivityLog, 1000),
	}
	go logger.processLogs()
	return logger
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		if al.rateLimiter.Allow(userID) {
			activity := ActivityLog{
				UserID:    userID,
				Action:    "request",
				Path:      r.URL.Path,
				Method:    r.Method,
				Timestamp: time.Now(),
				IPAddress: ip,
			}
			select {
			case al.logChannel <- activity:
			default:
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) processLogs() {
	for activity := range al.logChannel {
		jsonData, err := json.Marshal(activity)
		if err == nil {
			go func(data []byte) {
				_ = data
			}(jsonData)
		}
	}
}

func (al *ActivityLogger) Shutdown() {
	close(al.logChannel)
}