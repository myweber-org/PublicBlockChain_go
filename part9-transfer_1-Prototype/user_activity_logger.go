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
		"Activity: %s %s | IP: %s | Duration: %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package main

import (
	"fmt"
	"log"
	"os"
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
	UserID     string
	Action     ActivityType
	Timestamp  time.Time
	Additional map[string]interface{}
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

func (al *ActivityLogger) LogActivity(activity UserActivity) error {
	logEntry := fmt.Sprintf("[%s] User: %s, Action: %s",
		activity.Timestamp.Format(time.RFC3339),
		activity.UserID,
		activity.Action)

	if len(activity.Additional) > 0 {
		logEntry += fmt.Sprintf(", Details: %v", activity.Additional)
	}

	_, err := al.logFile.WriteString(logEntry + "\n")
	return err
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("user_activities.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	activities := []UserActivity{
		{
			UserID:    "user123",
			Action:    Login,
			Timestamp: time.Now(),
		},
		{
			UserID:    "user456",
			Action:    Purchase,
			Timestamp: time.Now(),
			Additional: map[string]interface{}{
				"item_id":  "prod789",
				"amount":   49.99,
				"currency": "USD",
			},
		},
		{
			UserID:    "user123",
			Action:    View,
			Timestamp: time.Now().Add(5 * time.Minute),
			Additional: map[string]interface{}{
				"page":     "/products/abc",
				"duration": 120,
			},
		},
	}

	for _, activity := range activities {
		if err := logger.LogActivity(activity); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}

	fmt.Println("Activities logged successfully")
}