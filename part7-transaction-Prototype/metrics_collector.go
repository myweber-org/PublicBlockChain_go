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
}