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

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    _, err = l.logFile.Write(append(data, '\n'))
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

    err = logger.LogActivity("user123", "login", "User logged in from IP 192.168.1.100")
    if err != nil {
        log.Println("Failed to log activity:", err)
    }

    err = logger.LogActivity("user123", "file_upload", "Uploaded document.pdf")
    if err != nil {
        log.Println("Failed to log activity:", err)
    }

    fmt.Println("Activity logging completed")
}