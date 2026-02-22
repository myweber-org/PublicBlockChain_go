package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    CPUPercent    float64
    MemoryAlloc   uint64
    MemoryTotal   uint64
    GoroutineCount int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:     time.Now(),
        MemoryAlloc:   m.Alloc,
        MemoryTotal:   m.TotalAlloc,
        GoroutineCount: runtime.NumGoroutine(),
        CPUPercent:    calculateCPUUsage(),
    }
}

func calculateCPUUsage() float64 {
    start := time.Now()
    startGoroutines := runtime.NumGoroutine()

    time.Sleep(100 * time.Millisecond)

    elapsed := time.Since(start)
    endGoroutines := runtime.NumGoroutine()

    usage := float64(endGoroutines-startGoroutines) / float64(runtime.NumCPU())
    if usage < 0 {
        usage = 0
    }
    return usage * 100 / elapsed.Seconds()
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
    fmt.Printf("Memory Allocated: %v bytes\n", metrics.MemoryAlloc)
    fmt.Printf("Total Memory: %v bytes\n", metrics.MemoryTotal)
    fmt.Printf("Active Goroutines: %d\n", metrics.GoroutineCount)
    fmt.Println("---")
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
}package main

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

    httpRequestTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(httpRequestTotal)
}

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

    handler := metricsMiddleware(mux)
    http.ListenAndServe(":8080", handler)
}