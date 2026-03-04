package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s | UserAgent: %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
			r.RemoteAddr,
			r.UserAgent(),
		)
	})
}package middleware

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
	
	log.Printf(
		"Method: %s | Path: %s | Duration: %v | Timestamp: %s",
		r.Method,
		r.URL.Path,
		duration,
		time.Now().Format(time.RFC3339),
	)
}package middleware

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
		startTime := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(startTime)
		
		al.Logger.Printf(
			"Method: %s | Path: %s | Status: %d | Duration: %v | User-Agent: %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
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
	startTime := time.Now()
	
	al.handler.ServeHTTP(w, r)
	
	duration := time.Since(startTime)
	
	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
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
	Method    string      `json:"method"`
	Path      string      `json:"path"`
	Status    int         `json:"status"`
	Duration  float64     `json:"duration_ms"`
	ClientIP  string      `json:"client_ip"`
	UserAgent string      `json:"user_agent"`
	Extra     interface{} `json:"extra,omitempty"`
}

type LoggerConfig struct {
	OutputFormat string
	IncludeExtra bool
	LogLevel     string
}

func ActivityLogger(config LoggerConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			recorder := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			
			next.ServeHTTP(recorder, r)
			
			duration := time.Since(start).Seconds() * 1000
			
			clientIP := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				clientIP = forwarded
			}
			
			logEntry := ActivityLog{
				Timestamp: time.Now().UTC(),
				Method:    r.Method,
				Path:      r.URL.Path,
				Status:    recorder.statusCode,
				Duration:  duration,
				ClientIP:  clientIP,
				UserAgent: r.UserAgent(),
			}
			
			if config.IncludeExtra {
				logEntry.Extra = map[string]interface{}{
					"protocol": r.Proto,
					"host":     r.Host,
					"referer":  r.Referer(),
				}
			}
			
			switch config.OutputFormat {
			case "json":
				data, err := json.Marshal(logEntry)
				if err != nil {
					log.Printf("Failed to marshal log entry: %v", err)
				} else {
					log.Println(string(data))
				}
			case "text":
				log.Printf("%s %s %d %.2fms %s %s",
					logEntry.Method,
					logEntry.Path,
					logEntry.Status,
					logEntry.Duration,
					logEntry.ClientIP,
					logEntry.UserAgent)
			default:
				log.Printf("Unsupported output format: %s", config.OutputFormat)
			}
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}