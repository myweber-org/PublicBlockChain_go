package metrics

import (
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
	
	i := 0
	for i < len(a.timestamps) && a.timestamps[i].Before(cutoff) {
		i++
	}
	
	if i > 0 {
		a.dataPoints = a.dataPoints[i:]
		a.timestamps = a.timestamps[i:]
	}
}

func (a *Aggregator) Average() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	if len(a.dataPoints) == 0 {
		return 0.0
	}
	
	var sum float64
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

func (a *Aggregator) Max() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	if len(a.dataPoints) == 0 {
		return 0.0
	}
	
	max := a.dataPoints[0]
	for _, v := range a.dataPoints[1:] {
		if v > max {
			max = v
		}
	}
	return max
}