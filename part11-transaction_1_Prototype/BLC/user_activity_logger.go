package main

import (
    "encoding/json"
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
    UserID    string       `json:"user_id"`
    Action    ActivityType `json:"action"`
    Timestamp time.Time    `json:"timestamp"`
    Details   string       `json:"details,omitempty"`
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

func (l *ActivityLogger) LogActivity(userID string, action ActivityType, details string) error {
    activity := UserActivity{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }
    return l.encoder.Encode(activity)
}

func (l *ActivityLogger) Close() error {
    return l.logFile.Close()
}

func main() {
    logger, err := NewActivityLogger("user_activities.jsonl")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    activities := []struct {
        userID string
        action ActivityType
        details string
    }{
        {"user_123", Login, "Successful authentication"},
        {"user_123", View, "Viewed product catalog"},
        {"user_456", Purchase, "Order #7890 completed"},
        {"user_123", Logout, "Session terminated"},
    }

    for _, act := range activities {
        if err := logger.LogActivity(act.userID, act.action, act.details); err != nil {
            log.Printf("Failed to log activity: %v", err)
        }
    }

    fmt.Println("Activity logging completed")
}