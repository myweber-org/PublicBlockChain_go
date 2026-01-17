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

type ActivityLogger struct {
	logFile *os.File
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ActivityLogger{logFile: file}, nil
}

func (al *ActivityLogger) LogActivity(userID, action, details string) error {
	logEntry := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	entryJSON, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	entryJSON = append(entryJSON, '\n')
	_, err = al.logFile.Write(entryJSON)
	return err
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	activities := []struct {
		userID  string
		action  string
		details string
	}{
		{"user_001", "LOGIN", "User logged in from IP 192.168.1.100"},
		{"user_001", "VIEW_PAGE", "Accessed dashboard page"},
		{"user_002", "REGISTER", "New user registration completed"},
		{"user_001", "LOGOUT", "User session terminated"},
	}

	for _, activity := range activities {
		err := logger.LogActivity(activity.userID, activity.action, activity.details)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}

	fmt.Println("Activity logging completed")
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	Logger *log.Logger
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{Logger: logger}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		al.Logger.Printf(
			"method=%s path=%s status=%d duration=%s remote=%s user_agent=%s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.RemoteAddr,
			r.UserAgent(),
		)
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