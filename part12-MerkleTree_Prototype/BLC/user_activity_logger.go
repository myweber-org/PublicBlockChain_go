package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	Logger *log.Logger
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{Logger: logger}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		al.Logger.Printf(
			"Method=%s Path=%s Status=%d Duration=%s RemoteAddr=%s UserAgent=%s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.RemoteAddr,
			r.UserAgent(),
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time   `json:"timestamp"`
	UserID    string      `json:"user_id"`
	Action    string      `json:"action"`
	Path      string      `json:"path"`
	Method    string      `json:"method"`
	IPAddress string      `json:"ip_address"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

type LoggerConfig struct {
	OutputFormat string
	LogToFile    bool
	FilePath     string
}

func ActivityLogger(config LoggerConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			userID := extractUserID(r)
			action := determineAction(r)
			
			lw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(lw, r)
			
			duration := time.Since(start)
			
			activity := ActivityLog{
				Timestamp: time.Now(),
				UserID:    userID,
				Action:    action,
				Path:      r.URL.Path,
				Method:    r.Method,
				IPAddress: r.RemoteAddr,
				Metadata: map[string]interface{}{
					"duration_ms": duration.Milliseconds(),
					"status_code": lw.statusCode,
					"user_agent":  r.UserAgent(),
				},
			}
			
			logActivity(activity, config)
		})
	}
}

func extractUserID(r *http.Request) string {
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		return "authenticated_user"
	}
	return "anonymous"
}

func determineAction(r *http.Request) string {
	switch r.Method {
	case http.MethodGet:
		return "view"
	case http.MethodPost:
		return "create"
	case http.MethodPut:
		return "update"
	case http.MethodDelete:
		return "delete"
	default:
		return "other"
	}
}

func logActivity(activity ActivityLog, config LoggerConfig) {
	switch config.OutputFormat {
	case "json":
		data, err := json.Marshal(activity)
		if err != nil {
			log.Printf("Failed to marshal activity log: %v", err)
			return
		}
		log.Println(string(data))
	case "text":
		log.Printf("[%s] %s %s %s %s (%dms)",
			activity.Timestamp.Format(time.RFC3339),
			activity.UserID,
			activity.Action,
			activity.Method,
			activity.Path,
			activity.Metadata.(map[string]interface{})["duration_ms"])
	default:
		log.Printf("Unsupported output format: %s", config.OutputFormat)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}