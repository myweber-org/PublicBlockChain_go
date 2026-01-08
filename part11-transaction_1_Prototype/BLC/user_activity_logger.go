package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityRecorder struct {
	ResponseWriter http.ResponseWriter
	StatusCode     int
}

func (ar *ActivityRecorder) WriteHeader(code int) {
	ar.StatusCode = code
	ar.ResponseWriter.WriteHeader(code)
}

func (ar *ActivityRecorder) Header() http.Header {
	return ar.ResponseWriter.Header()
}

func (ar *ActivityRecorder) Write(b []byte) (int, error) {
	return ar.ResponseWriter.Write(b)
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recorder := &ActivityRecorder{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(recorder, r)

		duration := time.Since(startTime)
		log.Printf(
			"Activity: %s %s %d %v %s",
			r.Method,
			r.URL.Path,
			recorder.StatusCode,
			duration,
			r.RemoteAddr,
		)
	})
}