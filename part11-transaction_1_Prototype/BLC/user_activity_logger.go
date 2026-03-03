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

func logActivity(userID, action, details string) {
	logEntry := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	logData, err := json.MarshalIndent(logEntry, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	logFile, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer logFile.Close()

	if _, err := logFile.Write(append(logData, '\n')); err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}
}

func main() {
	logActivity("user123", "LOGIN", "User logged in from IP 192.168.1.100")
	logActivity("user456", "FILE_UPLOAD", "Uploaded document.pdf (2.4 MB)")
	logActivity("user123", "LOGOUT", "Session duration: 1h 23m")

	fmt.Println("Activity logging completed. Check activity.log file.")
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
	log.Printf(
		"Method: %s | Path: %s | Status: %d | Duration: %v | UserAgent: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.UserAgent(),
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}