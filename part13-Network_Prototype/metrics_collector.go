
package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Timestamp   time.Time
	CPUPercent  float64
	MemoryAlloc uint64
	NumGoroutine int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:   time.Now(),
		CPUPercent:  getCPUUsage(),
		MemoryAlloc: m.Alloc,
		NumGoroutine: runtime.NumGoroutine(),
	}
}

func getCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)
	return float64(elapsed) / float64(time.Second) * 100
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics := collectMetrics()
		fmt.Printf("[%s] CPU: %.2f%% | Memory: %v bytes | Goroutines: %d\n",
			metrics.Timestamp.Format("15:04:05"),
			metrics.CPUPercent,
			metrics.MemoryAlloc,
			metrics.NumGoroutine)
	}
}