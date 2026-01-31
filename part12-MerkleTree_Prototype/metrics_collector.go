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
	time.Sleep(50 * time.Millisecond)
	elapsed := time.Since(start).Seconds()
	return (50.0 / 1000.0) / elapsed * 100
}

func displayMetrics(metrics SystemMetrics) {
	fmt.Printf("[%s] CPU: %.2f%% | Memory: %v MB | Goroutines: %d\n",
		metrics.Timestamp.Format("15:04:05"),
		metrics.CPUPercent,
		metrics.MemoryAlloc/1024/1024,
		metrics.NumGoroutine)
}

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics := collectMetrics()
		displayMetrics(metrics)
	}
}