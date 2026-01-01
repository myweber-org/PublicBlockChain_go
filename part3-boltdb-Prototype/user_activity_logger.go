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
	log.Printf(
		"%s %s %d %s %s",
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
	log.Printf(
		"%s %s %d %s %s",
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
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Endpoint  string
	Method    string
	Timestamp time.Time
	IPAddress string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserID(authHeader)
		}

		activity := ActivityLog{
			UserID:    userID,
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Timestamp: start,
			IPAddress: r.RemoteAddr,
		}

		log.Printf("Activity: %s %s by %s from %s", 
			activity.Method, 
			activity.Endpoint, 
			activity.UserID, 
			activity.IPAddress)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(token string) string {
	return "user_" + token[:8]
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

	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
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
	
	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}