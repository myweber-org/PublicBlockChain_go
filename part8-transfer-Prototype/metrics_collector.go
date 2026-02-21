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
}package main

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path", "status"},
    )

    httpRequestCount = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_request_count_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(rw, r)

        duration := time.Since(start).Seconds()
        status := http.StatusText(rw.statusCode)

        httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
        httpRequestCount.WithLabelValues(r.Method, r.URL.Path, status).Inc()
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, metrics!"))
    })
    mux.Handle("/metrics", promhttp.Handler())

    wrappedMux := metricsMiddleware(mux)
    http.ListenAndServe(":8080", wrappedMux)
}