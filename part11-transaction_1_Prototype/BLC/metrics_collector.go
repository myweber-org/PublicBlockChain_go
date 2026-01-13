
package main

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
	GoroutineCount int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:     time.Now(),
		CPUUsage:      calculateCPUUsage(),
		MemoryAlloc:   m.Alloc,
		MemoryTotal:   m.Sys,
		GoroutineCount: runtime.NumGoroutine(),
	}
}

func calculateCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(50 * time.Millisecond)
	elapsed := time.Since(start)

	return float64(elapsed) / float64(time.Second) * 100
}

func printMetrics(metrics SystemMetrics) {
	fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUUsage)
	fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Total Memory: %d bytes\n", metrics.MemoryTotal)
	fmt.Printf("Goroutines: %d\n", metrics.GoroutineCount)
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
}