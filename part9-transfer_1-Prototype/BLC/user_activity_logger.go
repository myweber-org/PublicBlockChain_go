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
}