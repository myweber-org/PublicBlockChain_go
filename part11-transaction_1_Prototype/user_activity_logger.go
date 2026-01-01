package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Details   string    `json:"details"`
}

func logActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(logEntry); err != nil {
        return err
    }

    return nil
}

func main() {
    if err := logActivity("user123", "LOGIN", "User logged in from IP 192.168.1.100"); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Activity logged successfully")
}package middleware

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
	userID := extractUserID(r)
	ipAddress := r.RemoteAddr

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("User %s from %s accessed %s %s - Duration: %v", userID, ipAddress, r.Method, r.URL.Path, duration)
}

func extractUserID(r *http.Request) string {
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return userID
	}
	return "anonymous"
}package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu      sync.RWMutex
	clients map[string]*clientActivity
	limit   int
	window  time.Duration
}

type clientActivity struct {
	count    int
	lastSeen time.Time
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		clients: make(map[string]*clientActivity),
		limit:   limit,
		window:  window,
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		path := r.URL.Path

		if !al.allowRequest(clientIP) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		log.Printf("Activity: %s - %s - %s", clientIP, userAgent, path)
		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) allowRequest(clientIP string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	activity, exists := al.clients[clientIP]

	if !exists {
		al.clients[clientIP] = &clientActivity{
			count:    1,
			lastSeen: now,
		}
		return true
	}

	if now.Sub(activity.lastSeen) > al.window {
		activity.count = 1
		activity.lastSeen = now
		return true
	}

	if activity.count >= al.limit {
		return false
	}

	activity.count++
	activity.lastSeen = now
	return true
}

func (al *ActivityLogger) Cleanup() {
	ticker := time.NewTicker(al.window * 2)
	go func() {
		for range ticker.C {
			al.mu.Lock()
			now := time.Now()
			for ip, activity := range al.clients {
				if now.Sub(activity.lastSeen) > al.window*2 {
					delete(al.clients, ip)
				}
			}
			al.mu.Unlock()
		}
	}()
}package middleware

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

	log.Printf("Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}