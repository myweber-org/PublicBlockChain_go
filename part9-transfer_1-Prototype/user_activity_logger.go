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
	writer := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	
	al.handler.ServeHTTP(writer, r)
	
	duration := time.Since(start)
	
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		writer.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}package main

import (
    "encoding/json"
    "fmt"
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
    logFile string
}

func NewActivityLogger(logFile string) *ActivityLogger {
    return &ActivityLogger{logFile: logFile}
}

func (al *ActivityLogger) LogActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    file, err := os.OpenFile(al.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(logEntry); err != nil {
        return fmt.Errorf("failed to encode log entry: %w", err)
    }

    return nil
}

func main() {
    logger := NewActivityLogger("user_activity.log")

    activities := []struct {
        userID  string
        action  string
        details string
    }{
        {"user_001", "LOGIN", "User logged in from IP 192.168.1.100"},
        {"user_001", "VIEW_PAGE", "Accessed dashboard page"},
        {"user_002", "REGISTER", "New user registration completed"},
        {"user_001", "UPDATE_PROFILE", "Changed email address"},
        {"user_003", "LOGOUT", "User session terminated"},
    }

    for _, activity := range activities {
        if err := logger.LogActivity(activity.userID, activity.action, activity.details); err != nil {
            fmt.Printf("Failed to log activity: %v\n", err)
        } else {
            fmt.Printf("Logged: %s - %s\n", activity.userID, activity.action)
        }
    }

    fmt.Println("Activity logging completed.")
}