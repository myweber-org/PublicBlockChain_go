
package main

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

    httpRequestTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
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
        httpRequestTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
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
    "log"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    GoroutineCount int
    MemoryAlloc   uint64
    MemoryTotal   uint64
    CPUCores      int
}

func collectMetrics() SystemMetrics {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    return SystemMetrics{
        Timestamp:     time.Now(),
        GoroutineCount: runtime.NumGoroutine(),
        MemoryAlloc:   memStats.Alloc,
        MemoryTotal:   memStats.TotalAlloc,
        CPUCores:      runtime.NumCPU(),
    }
}

func logMetrics(metrics SystemMetrics) {
    log.Printf(
        "Metrics collected at %s: Goroutines=%d, MemoryAlloc=%d bytes, TotalMemory=%d bytes, CPU Cores=%d",
        metrics.Timestamp.Format(time.RFC3339),
        metrics.GoroutineCount,
        metrics.MemoryAlloc,
        metrics.MemoryTotal,
        metrics.CPUCores,
    )
}

func main() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    fmt.Println("Starting system metrics collector...")

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            logMetrics(metrics)
        }
    }
}