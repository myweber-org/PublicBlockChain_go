
package metrics

import (
	"sync"
	"time"
)

type RequestMetrics struct {
	mu               sync.RWMutex
	latencyHistogram map[string][]time.Duration
	statusCounts     map[string]map[int]int64
	totalRequests    int64
}

func NewRequestMetrics() *RequestMetrics {
	return &RequestMetrics{
		latencyHistogram: make(map[string][]time.Duration),
		statusCounts:     make(map[string]map[int]int64),
	}
}

func (rm *RequestMetrics) RecordRequest(endpoint string, latency time.Duration, statusCode int) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.totalRequests++

	if _, exists := rm.latencyHistogram[endpoint]; !exists {
		rm.latencyHistogram[endpoint] = make([]time.Duration, 0, 1000)
	}
	rm.latencyHistogram[endpoint] = append(rm.latencyHistogram[endpoint], latency)

	if _, exists := rm.statusCounts[endpoint]; !exists {
		rm.statusCounts[endpoint] = make(map[int]int64)
	}
	rm.statusCounts[endpoint][statusCode]++
}

func (rm *RequestMetrics) GetEndpointStats(endpoint string) (avgLatency time.Duration, statusDistribution map[int]int64) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	latencies := rm.latencyHistogram[endpoint]
	if len(latencies) == 0 {
		return 0, nil
	}

	var total time.Duration
	for _, lat := range latencies {
		total += lat
	}
	avgLatency = total / time.Duration(len(latencies))

	statusDistribution = make(map[int]int64)
	for code, count := range rm.statusCounts[endpoint] {
		statusDistribution[code] = count
	}

	return avgLatency, statusDistribution
}

func (rm *RequestMetrics) GetTotalRequests() int64 {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.totalRequests
}