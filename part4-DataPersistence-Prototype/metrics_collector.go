package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    MemoryAlloc   uint64
    TotalAlloc    uint64
    SysMemory     uint64
    NumGoroutines int
    NumCPU        int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:     time.Now().UTC(),
        MemoryAlloc:   m.Alloc,
        TotalAlloc:    m.TotalAlloc,
        SysMemory:     m.Sys,
        NumGoroutines: runtime.NumGoroutine(),
        NumCPU:        runtime.NumCPU(),
    }
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Metrics collected at: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
    fmt.Printf("Total Allocated: %d bytes\n", metrics.TotalAlloc)
    fmt.Printf("System Memory: %d bytes\n", metrics.SysMemory)
    fmt.Printf("Goroutines: %d\n", metrics.NumGoroutines)
    fmt.Printf("CPU Cores: %d\n", metrics.NumCPU)
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