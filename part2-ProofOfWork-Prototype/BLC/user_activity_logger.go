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

	log.Printf("Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package main

import (
	"encoding/json"
	"fmt"
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
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ActivityLogger{logFile: file}, nil
}

func (l *ActivityLogger) LogActivity(userID, eventType, details string) error {
	event := ActivityEvent{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		EventType: eventType,
		Details:   details,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = l.logFile.Write(data)
	return err
}

func (l *ActivityLogger) Close() error {
	return l.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("activity.log")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	events := []struct {
		userID    string
		eventType string
		details   string
	}{
		{"user123", "login", "Successful authentication"},
		{"user123", "view_page", "Accessed dashboard"},
		{"user456", "purchase", "Order ID: 78910"},
		{"user123", "logout", "Session terminated"},
	}

	for _, e := range events {
		err := logger.LogActivity(e.userID, e.eventType, e.details)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}

	fmt.Println("Activity logging completed")
}