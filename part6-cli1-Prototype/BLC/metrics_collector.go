package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    MemoryAlloc uint64
    CPUCores    int
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:   time.Now().UTC(),
        MemoryAlloc: m.Alloc,
        CPUCores:    runtime.NumCPU(),
        Goroutines:  runtime.NumGoroutine(),
    }
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Time: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
    fmt.Printf("CPU Cores: %d\n", metrics.CPUCores)
    fmt.Printf("Active Goroutines: %d\n", metrics.Goroutines)
    fmt.Println("---")
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for i := 0; i < 3; i++ {
        metrics := collectMetrics()
        printMetrics(metrics)
        <-ticker.C
    }
}