
package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize  time.Duration
	dataPoints  []float64
	timestamps  []time.Time
	mu          sync.RWMutex
}

func NewAggregator(windowSize time.Duration) *Aggregator {
	return &Aggregator{
		windowSize: windowSize,
		dataPoints: make([]float64, 0),
		timestamps: make([]time.Time, 0),
	}
}

func (a *Aggregator) Add(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	now := time.Now()
	a.dataPoints = append(a.dataPoints, value)
	a.timestamps = append(a.timestamps, now)
	
	a.cleanup(now)
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
		a.dataPoints = a.dataPoints[validStart:]
		a.timestamps = a.timestamps[validStart:]
	}
}

func (a *Aggregator) Percentile(p float64) float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	if len(a.dataPoints) == 0 {
		return 0
	}
	
	sorted := make([]float64, len(a.dataPoints))
	copy(sorted, a.dataPoints)
	sort.Float64s(sorted)
	
	index := p * float64(len(sorted)-1) / 100.0
	lower := int(index)
	upper := lower + 1
	
	if upper >= len(sorted) {
		return sorted[lower]
	}
	
	weight := index - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

func (a *Aggregator) Average() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	if len(a.dataPoints) == 0 {
		return 0
	}
	
	sum := 0.0
	for _, v := range a.dataPoints {
		sum += v
	}
	return sum / float64(len(a.dataPoints))
}

func (a *Aggregator) Count() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.dataPoints)
}