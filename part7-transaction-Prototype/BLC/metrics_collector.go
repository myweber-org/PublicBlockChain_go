package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUPercent  float64
    MemoryUsage uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    return SystemMetrics{
        Timestamp:   time.Now(),
        CPUPercent:  getCPUUsage(),
        MemoryUsage: memStats.Alloc,
        Goroutines:  runtime.NumGoroutine(),
    }
}

func getCPUUsage() float64 {
    start := time.Now()
    startCPU := runtime.NumCPU()

    time.Sleep(100 * time.Millisecond)

    elapsed := time.Since(start)
    endCPU := runtime.NumCPU()

    if elapsed == 0 {
        return 0.0
    }
    return float64(endCPU-startCPU) / float64(runtime.NumCPU()) * 100.0
}

func displayMetrics(metrics SystemMetrics) {
    fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
    fmt.Printf("Memory Usage: %d bytes\n", metrics.MemoryUsage)
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