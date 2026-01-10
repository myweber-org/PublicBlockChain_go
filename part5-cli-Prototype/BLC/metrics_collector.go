package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Timestamp    time.Time
	CPUPercent   float64
	MemoryAlloc  uint64
	MemoryTotal  uint64
	GoroutineCnt int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:    time.Now(),
		CPUPercent:   getCPUUsage(),
		MemoryAlloc:  m.Alloc,
		MemoryTotal:  m.TotalAlloc,
		GoroutineCnt: runtime.NumGoroutine(),
	}
}

func getCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start).Seconds()
	return (100.0 - (elapsed * 1000 / 100)) / 100.0
}

func printMetrics(metrics SystemMetrics) {
	fmt.Printf("[%s] CPU: %.2f%% | Memory: %v/%v bytes | Goroutines: %d\n",
		metrics.Timestamp.Format("15:04:05"),
		metrics.CPUPercent,
		metrics.MemoryAlloc,
		metrics.MemoryTotal,
		metrics.GoroutineCnt)
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