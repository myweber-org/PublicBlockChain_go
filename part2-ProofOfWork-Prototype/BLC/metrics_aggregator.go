
package aggregator

import (
	"container/list"
	"sort"
	"sync"
	"time"
)

type MetricPoint struct {
	Value     float64
	Timestamp time.Time
}

type SlidingWindowAggregator struct {
	windowSize  time.Duration
	maxPoints   int
	points      *list.List
	mu          sync.RWMutex
	percentiles []float64
}

func NewSlidingWindowAggregator(windowSize time.Duration, maxPoints int, percentiles []float64) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		maxPoints:   maxPoints,
		points:      list.New(),
		percentiles: percentiles,
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.points.PushBack(&MetricPoint{
		Value:     value,
		Timestamp: now,
	})

	swa.cleanupOldPoints(now)
	if swa.points.Len() > swa.maxPoints {
		swa.points.Remove(swa.points.Front())
	}
}

func (swa *SlidingWindowAggregator) cleanupOldPoints(now time.Time) {
	cutoff := now.Add(-swa.windowSize)
	for e := swa.points.Front(); e != nil; {
		next := e.Next()
		if mp := e.Value.(*MetricPoint); mp.Timestamp.Before(cutoff) {
			swa.points.Remove(e)
		}
		e = next
	}
}

func (swa *SlidingWindowAggregator) GetAggregatedMetrics() map[string]float64 {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if swa.points.Len() == 0 {
		return make(map[string]float64)
	}

	values := make([]float64, 0, swa.points.Len())
	sum := 0.0
	min := 0.0
	max := 0.0
	first := true

	for e := swa.points.Front(); e != nil; e = e.Next() {
		value := e.Value.(*MetricPoint).Value
		values = append(values, value)
		sum += value

		if first {
			min = value
			max = value
			first = false
		} else {
			if value < min {
				min = value
			}
			if value > max {
				max = value
			}
		}
	}

	sort.Float64s(values)
	results := make(map[string]float64)
	results["count"] = float64(len(values))
	results["sum"] = sum
	results["avg"] = sum / float64(len(values))
	results["min"] = min
	results["max"] = max

	for _, p := range swa.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		idx := int(float64(len(values)-1) * p / 100.0)
		results[formatPercentileKey(p)] = values[idx]
	}

	return results
}

func formatPercentileKey(p float64) string {
	return "p" + strings.Replace(fmt.Sprintf("%.1f", p), ".", "_", -1)
}