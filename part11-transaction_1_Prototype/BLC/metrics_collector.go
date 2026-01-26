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