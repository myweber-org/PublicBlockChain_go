package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

type Metrics struct {
    RequestCount    int
    TotalLatency    time.Duration
    ErrorCount      int
}

var metrics = Metrics{}

func handler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    defer func() {
        latency := time.Since(start)
        metrics.RequestCount++
        metrics.TotalLatency += latency
        
        if r.URL.Path == "/error" {
            metrics.ErrorCount++
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "Error occurred")
            return
        }
        
        fmt.Fprintf(w, "Request processed in %v", latency)
    }()
    
    time.Sleep(10 * time.Millisecond)
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
    fmt.Fprintf(w, "Average Latency: %v\n", avgLatency)
    fmt.Fprintf(w, "Error Rate: %.2f%%\n", errorRate)
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/metrics", metricsHandler)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}