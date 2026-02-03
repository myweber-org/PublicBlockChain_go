
package aggregator

import (
	"sort"
	"sync"
	"time"
)

type Metric struct {
	Value     float64
	Timestamp time.Time
}

type SlidingWindow struct {
	windowSize  time.Duration
	maxSamples  int
	metrics     []Metric
	mu          sync.RWMutex
}

func NewSlidingWindow(windowSize time.Duration, maxSamples int) *SlidingWindow {
	return &SlidingWindow{
		windowSize: windowSize,
		maxSamples: maxSamples,
		metrics:    make([]Metric, 0, maxSamples),
	}
}

func (sw *SlidingWindow) Add(value float64) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	sw.metrics = append(sw.metrics, Metric{Value: value, Timestamp: now})
	sw.cleanup()
}

func (sw *SlidingWindow) cleanup() {
	cutoff := time.Now().Add(-sw.windowSize)
	validStart := 0

	for i, metric := range sw.metrics {
		if metric.Timestamp.After(cutoff) {
			validStart = i
			break
		}
	}

	sw.metrics = sw.metrics[validStart:]

	if len(sw.metrics) > sw.maxSamples {
		sw.metrics = sw.metrics[len(sw.metrics)-sw.maxSamples:]
	}
}

func (sw *SlidingWindow) Percentile(p float64) float64 {
	sw.mu.RLock()
	defer sw.mu.RUnlock()

	if len(sw.metrics) == 0 {
		return 0
	}

	values := make([]float64, len(sw.metrics))
	for i, m := range sw.metrics {
		values[i] = m.Value
	}

	sort.Float64s(values)

	index := p * float64(len(values)-1)
	lower := int(index)
	upper := lower + 1
	weight := index - float64(lower)

	if upper >= len(values) {
		return values[lower]
	}

	return values[lower]*(1-weight) + values[upper]*weight
}

func (sw *SlidingWindow) Count() int {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	return len(sw.metrics)
}

func (sw *SlidingWindow) Average() float64 {
	sw.mu.RLock()
	defer sw.mu.RUnlock()

	if len(sw.metrics) == 0 {
		return 0
	}

	sum := 0.0
	for _, m := range sw.metrics {
		sum += m.Value
	}
	return sum / float64(len(sw.metrics))
}