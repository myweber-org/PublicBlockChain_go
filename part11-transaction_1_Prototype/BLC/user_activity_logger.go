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
}