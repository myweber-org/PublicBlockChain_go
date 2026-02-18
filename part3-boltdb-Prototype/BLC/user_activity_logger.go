package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/go-redis/redis/v8"
    "golang.org/x/time/rate"
)

type ActivityLogger struct {
    redisClient *redis.Client
    limiter     *rate.Limiter
}

func NewActivityLogger(redisAddr string, rps int) *ActivityLogger {
    rdb := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: "",
        DB:       0,
    })

    return &ActivityLogger{
        redisClient: rdb,
        limiter:     rate.NewLimiter(rate.Limit(rps), rps*2),
    }
}

func (al *ActivityLogger) LogActivity(userID, action string) error {
    if !al.limiter.Allow() {
        return fmt.Errorf("rate limit exceeded")
    }

    ctx := context.Background()
    key := fmt.Sprintf("activity:%s:%d", userID, time.Now().Unix())
    
    activity := map[string]interface{}{
        "user_id":    userID,
        "action":     action,
        "timestamp":  time.Now().Format(time.RFC3339),
        "user_agent": "web-client",
    }

    err := al.redisClient.HSet(ctx, key, activity).Err()
    if err != nil {
        return fmt.Errorf("failed to log activity: %w", err)
    }

    al.redisClient.Expire(ctx, key, 24*time.Hour)
    return nil
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Header.Get("X-User-ID")
        if userID != "" {
            go al.LogActivity(userID, r.URL.Path)
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    logger := NewActivityLogger("localhost:6379", 10)
    
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Request processed")
    })

    wrappedMux := logger.Middleware(mux)
    
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", wrappedMux)
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

	log.Printf(
		"Method: %s | Path: %s | Duration: %v | RemoteAddr: %s",
		r.Method,
		r.URL.Path,
		duration,
		r.RemoteAddr,
	)
}