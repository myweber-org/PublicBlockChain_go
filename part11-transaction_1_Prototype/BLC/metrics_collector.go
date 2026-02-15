package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUPercent  float64
    MemoryUsed  uint64
    MemoryTotal uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    metrics := SystemMetrics{
        Timestamp:   time.Now(),
        MemoryUsed:  m.Alloc,
        MemoryTotal: m.Sys,
        Goroutines:  runtime.NumGoroutine(),
    }

    // Simulate CPU usage calculation
    metrics.CPUPercent = calculateCPUUsage()
    return metrics
}

func calculateCPUUsage() float64 {
    // Placeholder for actual CPU calculation logic
    // In production, use gopsutil or similar library
    return 45.7 // Simulated value
}

func displayMetrics(metrics SystemMetrics) {
    fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
    fmt.Printf("Memory Used: %v MB\n", metrics.MemoryUsed/1024/1024)
    fmt.Printf("Memory Total: %v MB\n", metrics.MemoryTotal/1024/1024)
    fmt.Printf("Active Goroutines: %d\n", metrics.Goroutines)
    fmt.Println("---")
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            displayMetrics(metrics)
        }
    }
}
package main

import (
    "fmt"
    "net/http"
    "time"
)

type Metrics struct {
    RequestCount    int
    TotalLatency    time.Duration
    StatusCodes     map[int]int
}

var metrics = Metrics{
    StatusCodes: make(map[int]int),
}

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(recorder, r)
        
        duration := time.Since(start)
        metrics.RequestCount++
        metrics.TotalLatency += duration
        metrics.StatusCodes[recorder.statusCode]++
    })
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
    rr.statusCode = code
    rr.ResponseWriter.WriteHeader(code)
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
    mux := http.NewServeMux()
    mux.Handle("/", metricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(50 * time.Millisecond)
        w.Write([]byte("Hello, World!"))
    })))
    mux.HandleFunc("/metrics", metricsHandler)
    
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", mux)
}package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Timestamp    time.Time
	CPUUsage     float64
	MemoryAlloc  uint64
	MemoryTotal  uint64
	GoroutineCnt int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:    time.Now(),
		CPUUsage:     getCPUUsage(),
		MemoryAlloc:  m.Alloc,
		MemoryTotal:  m.Sys,
		GoroutineCnt: runtime.NumGoroutine(),
	}
}

func getCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(50 * time.Millisecond)
	elapsed := time.Since(start).Seconds()
	return (50.0 / 1000.0) / elapsed * 100.0
}

func printMetrics(metrics SystemMetrics) {
	fmt.Printf("Time: %v\n", metrics.Timestamp.Format("15:04:05"))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUUsage)
	fmt.Printf("Memory Allocated: %v bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Total Memory: %v bytes\n", metrics.MemoryTotal)
	fmt.Printf("Goroutines: %d\n", metrics.GoroutineCnt)
	fmt.Println("---")
}

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
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
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func instrumentedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		defer func() {
			duration := time.Since(start).Seconds()
			httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
			httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.statusCode)).Inc()
		}()

		next.ServeHTTP(rw, r)
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

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Application is running"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.Handle("/metrics", promhttp.Handler())

	instrumentedMux := instrumentedHandler(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      instrumentedMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}