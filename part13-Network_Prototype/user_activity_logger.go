package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	Format string
	Logger *log.Logger
}

func NewActivityLogger(format string, logger *log.Logger) *ActivityLogger {
	if format == "" {
		format = "default"
	}
	return &ActivityLogger{Format: format, Logger: logger}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		al.logEntry(r, recorder.statusCode, duration)
	})
}

func (al *ActivityLogger) logEntry(r *http.Request, status int, duration time.Duration) {
	switch al.Format {
	case "json":
		al.Logger.Printf(`{"time":"%s","method":"%s","path":"%s","status":%d,"duration_ms":%d}`,
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			status,
			duration.Milliseconds())
	case "simple":
		al.Logger.Printf("%s %s %d %v", r.Method, r.URL.Path, status, duration)
	default:
		al.Logger.Printf("[%s] %s %s %d %v",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			status,
			duration)
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}