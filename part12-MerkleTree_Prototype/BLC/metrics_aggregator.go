
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

func (swa *SlidingWindowAggregator) Add(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.metrics = append(swa.metrics, Metric{Timestamp: now, Value: value})

	cutoff := now.Add(-swa.windowSize)
	validStart := 0
	for i, m := range swa.metrics {
		if m.Timestamp.After(cutoff) {
			validStart = i
			break
		}
	}
	swa.metrics = swa.metrics[validStart:]

	if len(swa.metrics) > swa.maxSamples {
		swa.metrics = swa.metrics[len(swa.metrics)-swa.maxSamples:]
	}
}

func (swa *SlidingWindowAggregator) Percentile(p float64) (float64, bool) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if len(swa.metrics) == 0 {
		return 0, false
	}

	values := make([]float64, len(swa.metrics))
	for i, m := range swa.metrics {
		values[i] = m.Value
	}
	sort.Float64s(values)

	index := p * float64(len(values)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return values[lower], true
	}
	weight := index - float64(lower)
	return values[lower]*(1-weight) + values[upper]*weight, true
}

func (swa *SlidingWindowAggregator) Stats() (min, max, avg float64, count int) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if len(swa.metrics) == 0 {
		return 0, 0, 0, 0
	}

	min = swa.metrics[0].Value
	max = swa.metrics[0].Value
	sum := 0.0

	for _, m := range swa.metrics {
		if m.Value < min {
			min = m.Value
		}
		if m.Value > max {
			max = m.Value
		}
		sum += m.Value
	}

	return min, max, sum / float64(len(swa.metrics)), len(swa.metrics)
}

func main() {
	aggregator := NewSlidingWindowAggregator(5*time.Minute, 1000)

	for i := 0; i < 50; i++ {
		value := 50 + float64(i%20)
		aggregator.Add(value)
		time.Sleep(100 * time.Millisecond)
	}

	min, max, avg, count := aggregator.Stats()
	fmt.Printf("Samples: %d, Min: %.2f, Max: %.2f, Avg: %.2f\n", count, min, max, avg)

	if p95, ok := aggregator.Percentile(0.95); ok {
		fmt.Printf("95th percentile: %.2f\n", p95)
	}
}