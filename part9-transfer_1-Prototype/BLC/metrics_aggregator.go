package metrics

import (
	"sync"
	"time"
)

type Aggregator struct {
	windowSize  time.Duration
	metrics     []float64
	timestamps  []time.Time
	mu          sync.RWMutex
}

func NewAggregator(windowSize time.Duration) *Aggregator {
	return &Aggregator{
		windowSize: windowSize,
		metrics:    make([]float64, 0),
		timestamps: make([]time.Time, 0),
	}
}

func (a *Aggregator) AddMetric(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	a.metrics = append(a.metrics, value)
	a.timestamps = append(a.timestamps, now)

	a.cleanup(now)
}

func (a *Aggregator) GetAverage() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.metrics) == 0 {
		return 0
	}

	var sum float64
	for _, v := range a.metrics {
		sum += v
	}
	return sum / float64(len(a.metrics))
}

func (a *Aggregator) cleanup(currentTime time.Time) {
	cutoff := currentTime.Add(-a.windowSize)
	validStart := 0

	for i, ts := range a.timestamps {
		if ts.After(cutoff) {
			validStart = i
			break
		}
	}

	if validStart > 0 {
		a.metrics = a.metrics[validStart:]
		a.timestamps = a.timestamps[validStart:]
	}
}