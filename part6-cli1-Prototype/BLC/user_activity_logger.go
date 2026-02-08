
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
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(start)
	log.Printf("%s %s %d %v %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
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
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Status    int       `json:"status"`
	Duration  float64   `json:"duration_ms"`
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start).Seconds() * 1000
		
		logEntry := ActivityLog{
			Timestamp: time.Now().UTC(),
			Method:    r.Method,
			Path:      r.URL.Path,
			IP:        r.RemoteAddr,
			UserAgent: r.UserAgent(),
			Status:    recorder.statusCode,
			Duration:  duration,
		}
		
		logData, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("Failed to marshal activity log: %v", err)
			return
		}
		
		log.Println(string(logData))
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}