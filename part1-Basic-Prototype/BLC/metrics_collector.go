package main

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path", "status"},
    )

    httpRequestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(httpRequestCount)
}

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
        w.Write([]byte("Hello, World!"))
    })
    mux.Handle("/metrics", promhttp.Handler())

    wrappedMux := metricsMiddleware(mux)
    http.ListenAndServe(":8080", wrappedMux)
}package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    MemoryAlloc   uint64
    TotalAlloc    uint64
    Sys           uint64
    NumGC         uint32
    NumGoroutine  int
    CPUCount      int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:     time.Now().UTC(),
        MemoryAlloc:   m.Alloc,
        TotalAlloc:    m.TotalAlloc,
        Sys:           m.Sys,
        NumGC:         m.NumGC,
        NumGoroutine:  runtime.NumGoroutine(),
        CPUCount:      runtime.NumCPU(),
    }
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("=== System Metrics at %v ===\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("Goroutines: %d\n", metrics.NumGoroutine)
    fmt.Printf("CPU Cores: %d\n", metrics.CPUCount)
    fmt.Printf("Memory Allocated: %v MB\n", metrics.MemoryAlloc/1024/1024)
    fmt.Printf("Total Allocated: %v MB\n", metrics.TotalAlloc/1024/1024)
    fmt.Printf("System Memory: %v MB\n", metrics.Sys/1024/1024)
    fmt.Printf("Garbage Collections: %d\n", metrics.NumGC)
    fmt.Println()
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            printMetrics(metrics)
        }
    }
}