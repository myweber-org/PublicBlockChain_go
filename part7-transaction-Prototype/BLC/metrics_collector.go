package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time `json:"timestamp"`
    CPUUsage      float64   `json:"cpu_usage"`
    MemoryUsage   uint64    `json:"memory_usage"`
    GoroutineCount int      `json:"goroutine_count"`
    Hostname      string    `json:"hostname"`
}

func collectMetrics() SystemMetrics {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    hostname, err := os.Hostname()
    if err != nil {
        hostname = "unknown"
    }

    return SystemMetrics{
        Timestamp:     time.Now().UTC(),
        CPUUsage:      calculateCPUUsage(),
        MemoryUsage:   memStats.Alloc,
        GoroutineCount: runtime.NumGoroutine(),
        Hostname:      hostname,
    }
}

func calculateCPUUsage() float64 {
    start := time.Now()
    runtime.Gosched()
    time.Sleep(100 * time.Millisecond)
    elapsed := time.Since(start)
    return elapsed.Seconds() * 10
}

func main() {
    metrics := collectMetrics()

    jsonData, err := json.MarshalIndent(metrics, "", "  ")
    if err != nil {
        log.Fatal("Failed to marshal metrics:", err)
    }

    fmt.Println(string(jsonData))
}package main

import (
	"log"
	"net/http"
	"time"
)

var (
	requestLatency = make(map[string]time.Duration)
	errorCount     = make(map[string]int)
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		defer func() {
			duration := time.Since(start)
			requestLatency[path] = duration

			if r.Response != nil && r.Response.StatusCode >= 400 {
				errorCount[path]++
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	for path, latency := range requestLatency {
		errors := errorCount[path]
		log.Printf("Path: %s, Latency: %v, Errors: %d", path, latency, errors)
	}
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.Handle("/", metricsMiddleware(handler))
	http.Handle("/metrics", http.HandlerFunc(metricsHandler))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}