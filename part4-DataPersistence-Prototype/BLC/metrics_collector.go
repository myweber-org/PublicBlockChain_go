package main

import (
	"log"
	"net/http"
	"time"
)

var (
	requestCount    = make(map[string]int)
	requestDuration = make(map[string]time.Duration)
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		requestCount[path]++
		requestDuration[path] += duration

		log.Printf("Request to %s took %v, status: %d", path, duration, lrw.statusCode)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		for path, count := range requestCount {
			duration := requestDuration[path]
			avgDuration := duration / time.Duration(count)
			w.Write([]byte(`{"path":"` + path + `","count":` + string(count) + `,"avg_duration":` + avgDuration.String() + `}`))
		}
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", metricsMiddleware(mux)))
}