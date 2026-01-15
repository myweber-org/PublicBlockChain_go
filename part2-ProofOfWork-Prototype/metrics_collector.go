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
	startCPU := runtime.NumCgoCall()
	
	time.Sleep(100 * time.Millisecond)
	
	end := time.Now()
	endCPU := runtime.NumCgoCall()
	
	elapsed := end.Sub(start).Seconds()
	cpuDelta := float64(endCPU - startCPU)
	
	if elapsed > 0 {
		return cpuDelta / elapsed * 100
	}
	return 0.0
}

func displayMetrics(metrics SystemMetrics) {
	fmt.Printf("[%s] CPU: %.2f%% | Memory: %v bytes | Goroutines: %d\n",
		metrics.Timestamp.Format("15:04:05"),
		metrics.CPUPercent,
		metrics.MemoryAlloc,
		metrics.NumGoroutine)
}

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	fmt.Println("Starting system metrics collection...")
	
	for range ticker.C {
		metrics := collectMetrics()
		displayMetrics(metrics)
	}
}