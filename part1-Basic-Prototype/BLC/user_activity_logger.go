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
	Details   string    `json:"details,omitempty"`
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

func (l *ActivityLogger) LogActivity(userID, action, details string) error {
	entry := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}
	return l.encoder.Encode(entry)
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

	activities := []struct {
		userID  string
		action  string
		details string
	}{
		{"user_001", "LOGIN", "Successful authentication"},
		{"user_002", "UPLOAD", "File: report.pdf"},
		{"user_001", "LOGOUT", "Session ended"},
	}

	for _, act := range activities {
		if err := logger.LogActivity(act.userID, act.action, act.details); err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}

	fmt.Println("Activity logging completed")
}