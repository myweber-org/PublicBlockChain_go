package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	CPUUsage    float64
	MemoryAlloc uint64
	MemoryTotal uint64
	NumGoroutine int
	Timestamp   time.Time
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		CPUUsage:    getCPUUsage(),
		MemoryAlloc: m.Alloc,
		MemoryTotal: m.Sys,
		NumGoroutine: runtime.NumGoroutine(),
		Timestamp:   time.Now(),
	}
}

func getCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(50 * time.Millisecond)
	elapsed := time.Since(start)
	return float64(elapsed) / float64(time.Second) * 100
}

func displayMetrics(metrics SystemMetrics) {
	fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUUsage)
	fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Total Memory: %d bytes\n", metrics.MemoryTotal)
	fmt.Printf("Active Goroutines: %d\n", metrics.NumGoroutine)
	fmt.Println("---")
}

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
		select {
		case <-ticker.C:
			metrics := collectMetrics()
			displayMetrics(metrics)
		}
	}
}