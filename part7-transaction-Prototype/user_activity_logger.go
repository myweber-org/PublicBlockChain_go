package middleware

import (
	"context"
	"net/http"
	"time"
)

type ActivityLogger struct {
	store     ActivityStore
	rateLimit time.Duration
}

type ActivityStore interface {
	LogActivity(ctx context.Context, userID string, action string, metadata map[string]interface{}) error
}

func NewActivityLogger(store ActivityStore, rateLimit time.Duration) *ActivityLogger {
	return &ActivityLogger{
		store:     store,
		rateLimit: rateLimit,
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := extractUserID(r)
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}

		action := r.Method + " " + r.URL.Path
		metadata := map[string]interface{}{
			"user_agent": r.UserAgent(),
			"ip_address": r.RemoteAddr,
			"timestamp":  time.Now().UTC(),
		}

		go func() {
			select {
			case <-time.After(al.rateLimit):
				if err := al.store.LogActivity(ctx, userID, action, metadata); err != nil {
					logError(ctx, err)
				}
			case <-ctx.Done():
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return parseToken(auth)
	}
	return ""
}

func parseToken(token string) string {
	return token
}

func logError(ctx context.Context, err error) {
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
	writer := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	al.handler.ServeHTTP(writer, r)

	duration := time.Since(startTime)
	log.Printf(
		"Activity: %s %s | Status: %d | Duration: %v | User-Agent: %s",
		r.Method,
		r.URL.Path,
		writer.statusCode,
		duration,
		r.UserAgent(),
	)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}