package main

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

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	eventJSON = append(eventJSON, '\n')
	_, err = l.logFile.Write(eventJSON)
	return err
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
		{"user123", "login", "Successful authentication"},
		{"user456", "purchase", "Order #789 completed"},
		{"user123", "logout", "Session terminated"},
	}

	for _, e := range events {
		err := logger.LogActivity(e.userID, e.eventType, e.details)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		} else {
			fmt.Printf("Logged %s event for user %s\n", e.eventType, e.userID)
		}
	}

	fmt.Println("Activity logging completed")
}