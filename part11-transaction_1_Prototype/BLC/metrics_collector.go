
package main

import (
	"fmt"
	"net/http"
	"time"
)

type Metrics struct {
	RequestCount int
	TotalLatency time.Duration
	StatusCodes  map[int]int
}

var metrics = Metrics{
	StatusCodes: make(map[int]int),
}

func metricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		metrics.RequestCount++
		metrics.TotalLatency += duration
		metrics.StatusCodes[recorder.statusCode]++
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request processed")
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	avgLatency := time.Duration(0)
	if metrics.RequestCount > 0 {
		avgLatency = metrics.TotalLatency / time.Duration(metrics.RequestCount)
	}

	fmt.Fprintf(w, "Requests: %d\n", metrics.RequestCount)
	fmt.Fprintf(w, "Average Latency: %v\n", avgLatency)
	fmt.Fprintf(w, "Status Codes:\n")
	for code, count := range metrics.StatusCodes {
		fmt.Fprintf(w, "  %d: %d\n", code, count)
	}
}

func main() {
	http.HandleFunc("/", metricsMiddleware(handler))
	http.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe(":8080", nil)
}