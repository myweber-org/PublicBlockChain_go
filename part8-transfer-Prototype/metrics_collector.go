package main

import (
    "fmt"
    "net/http"
    "time"
)

type Metrics struct {
    RequestCount    int
    TotalLatency   time.Duration
    ErrorCount     int
}

var metrics = &Metrics{}

func metricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        recorder := &responseRecorder{w, 200}
        next(recorder, r)
        latency := time.Since(start)

        metrics.RequestCount++
        metrics.TotalLatency += latency
        if recorder.status >= 400 {
            metrics.ErrorCount++
        }
    }
}

type responseRecorder struct {
    http.ResponseWriter
    status int
}

func (r *responseRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

func handler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(10 * time.Millisecond)
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Request processed")
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
    avgLatency := time.Duration(0)
    if metrics.RequestCount > 0 {
        avgLatency = metrics.TotalLatency / time.Duration(metrics.RequestCount)
    }
    errorRate := 0.0
    if metrics.RequestCount > 0 {
        errorRate = float64(metrics.ErrorCount) / float64(metrics.RequestCount) * 100
    }

    fmt.Fprintf(w, "Requests: %d\n", metrics.RequestCount)
    fmt.Fprintf(w, "Avg Latency: %v\n", avgLatency)
    fmt.Fprintf(w, "Error Rate: %.2f%%\n", errorRate)
}

func main() {
    http.HandleFunc("/", metricsMiddleware(handler))
    http.HandleFunc("/metrics", metricsHandler)
    http.ListenAndServe(":8080", nil)
}