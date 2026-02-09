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
}