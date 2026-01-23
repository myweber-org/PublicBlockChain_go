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
	activity := ActivityLog{
		Timestamp: time.Now(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(activity); err != nil {
		log.Printf("Failed to encode activity: %v", err)
		return
	}

	fmt.Printf("Logged activity: %s performed %s\n", userID, action)
}

func main() {
	logActivity("user123", "login", "User logged in from web browser")
	logActivity("user456", "file_upload", "Uploaded document.pdf")
	logActivity("user123", "logout", "Session expired after 30 minutes")
}