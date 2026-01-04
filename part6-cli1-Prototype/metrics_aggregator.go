
package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

type Metric struct {
	Timestamp time.Time
	Value     float64
}

type SlidingWindowAggregator struct {
	windowSize  time.Duration
	maxSamples  int
	metrics     []Metric
	mu          sync.RWMutex
}

func NewSlidingWindowAggregator(windowSize time.Duration, maxSamples int) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize: windowSize,
		maxSamples: maxSamples,
		metrics:    make([]Metric, 0, maxSamples),
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.metrics = append(swa.metrics, Metric{Timestamp: now, Value: value})

	cutoff := now.Add(-swa.windowSize)
	validStart := 0
	for i, metric := range swa.metrics {
		if metric.Timestamp.After(cutoff) {
			validStart = i
			break
		}
	}
	swa.metrics = swa.metrics[validStart:]

	if len(swa.metrics) > swa.maxSamples {
		swa.metrics = swa.metrics[len(swa.metrics)-swa.maxSamples:]
	}
}

func (swa *SlidingWindowAggregator) CalculatePercentile(p float64) (float64, error) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if len(swa.metrics) == 0 {
		return 0, fmt.Errorf("no metrics available")
	}

	values := make([]float64, len(swa.metrics))
	for i, metric := range swa.metrics {
		values[i] = metric.Value
	}
	sort.Float64s(values)

	index := p * float64(len(values)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return values[lower], nil
	}
	weight := index - float64(lower)
	return values[lower]*(1-weight) + values[upper]*weight, nil
}

func (swa *SlidingWindowAggregator) GetStats() (float64, float64, float64, error) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if len(swa.metrics) == 0 {
		return 0, 0, 0, fmt.Errorf("no metrics available")
	}

	var sum, min, max float64
	min = math.MaxFloat64
	max = -math.MaxFloat64

	for _, metric := range swa.metrics {
		sum += metric.Value
		if metric.Value < min {
			min = metric.Value
		}
		if metric.Value > max {
			max = metric.Value
		}
	}

	avg := sum / float64(len(swa.metrics))
	return avg, min, max, nil
}

func main() {
	aggregator := NewSlidingWindowAggregator(5*time.Minute, 1000)

	for i := 0; i < 50; i++ {
		value := 10.0 + float64(i%20)*0.5
		aggregator.AddMetric(value)
		time.Sleep(100 * time.Millisecond)
	}

	avg, min, max, err := aggregator.GetStats()
	if err != nil {
		fmt.Printf("Error getting stats: %v\n", err)
		return
	}
	fmt.Printf("Average: %.2f, Min: %.2f, Max: %.2f\n", avg, min, max)

	p95, err := aggregator.CalculatePercentile(0.95)
	if err != nil {
		fmt.Printf("Error calculating percentile: %v\n", err)
		return
	}
	fmt.Printf("95th percentile: %.2f\n", p95)
}