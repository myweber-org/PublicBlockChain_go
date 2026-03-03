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
	Resource  string    `json:"resource"`
	Details   string    `json:"details"`
}

func NewActivityLog(userID, action, resource, details string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Details:   details,
	}
}

func (al *ActivityLog) ToJSON() ([]byte, error) {
	return json.MarshalIndent(al, "", "  ")
}

func LogActivity(logger *log.Logger, userID, action, resource, details string) {
	activity := NewActivityLog(userID, action, resource, details)
	jsonData, err := activity.ToJSON()
	if err != nil {
		logger.Printf("Failed to marshal activity log: %v", err)
		return
	}
	logger.Println(string(jsonData))
}

func main() {
	logFile, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	activityLogger := log.New(logFile, "", 0)

	LogActivity(activityLogger, "user123", "LOGIN", "auth", "User logged in from IP 192.168.1.100")
	LogActivity(activityLogger, "user456", "CREATE", "document", "Created new document 'Project Plan'")
	LogActivity(activityLogger, "user123", "UPDATE", "profile", "Changed email address")

	fmt.Println("Activity logging completed. Check activity.log file.")
}