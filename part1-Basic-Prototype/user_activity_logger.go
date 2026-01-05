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
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(start)
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}package main

import (
    "encoding/json"
    "log"
    "os"
    "sync"
    "time"
)

type ActivityType string

const (
    Login    ActivityType = "LOGIN"
    Logout   ActivityType = "LOGOUT"
    Purchase ActivityType = "PURCHASE"
    View     ActivityType = "VIEW"
)

type UserActivity struct {
    UserID    string       `json:"user_id"`
    Action    ActivityType `json:"action"`
    Timestamp time.Time    `json:"timestamp"`
    Metadata  interface{}  `json:"metadata,omitempty"`
}

type ActivityLogger struct {
    mu     sync.Mutex
    file   *os.File
    encoder *json.Encoder
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &ActivityLogger{
        file:    file,
        encoder: json.NewEncoder(file),
    }, nil
}

func (l *ActivityLogger) Log(activity UserActivity) error {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    activity.Timestamp = time.Now().UTC()
    return l.encoder.Encode(activity)
}

func (l *ActivityLogger) Close() error {
    return l.file.Close()
}

func main() {
    logger, err := NewActivityLogger("user_activities.jsonl")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    activities := []UserActivity{
        {UserID: "user_001", Action: Login, Metadata: map[string]string{"ip": "192.168.1.1"}},
        {UserID: "user_002", Action: View, Metadata: map[string]string{"page": "/products/123"}},
        {UserID: "user_001", Action: Purchase, Metadata: map[string]interface{}{"item_id": "prod_456", "amount": 29.99}},
    }

    for _, activity := range activities {
        if err := logger.Log(activity); err != nil {
            log.Printf("Failed to log activity: %v", err)
        }
    }
    
    log.Println("Activity logging completed")
}