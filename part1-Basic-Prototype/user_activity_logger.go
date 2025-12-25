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

	log.Printf(
		"Method: %s | Path: %s | Duration: %v | UserAgent: %s",
		r.Method,
		r.URL.Path,
		duration,
		r.UserAgent(),
	)
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
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	
	al.handler.ServeHTTP(recorder, r)
	
	duration := time.Since(start)
	log.Printf("[%s] %s %s - %d %v",
		r.RemoteAddr,
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
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
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityEvent struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Details   string    `json:"details"`
}

type ActivityLogger struct {
	logFile *os.File
	encoder *json.Encoder
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ActivityLogger{
		logFile: file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (l *ActivityLogger) LogActivity(userID, eventType, details string) error {
	event := ActivityEvent{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		EventType: eventType,
		Details:   details,
	}
	return l.encoder.Encode(event)
}

func (l *ActivityLogger) Close() error {
	return l.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	events := []struct {
		userID    string
		eventType string
		details   string
	}{
		{"user123", "LOGIN", "Successful authentication"},
		{"user123", "VIEW_PAGE", "/dashboard"},
		{"user456", "REGISTER", "New account created"},
		{"user123", "LOGOUT", "Session terminated"},
	}

	for _, e := range events {
		if err := logger.LogActivity(e.userID, e.eventType, e.details); err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}

	fmt.Println("Activity logging completed")
}