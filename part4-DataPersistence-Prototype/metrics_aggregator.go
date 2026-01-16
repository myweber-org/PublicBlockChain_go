
package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

type MetricPoint struct {
	Timestamp time.Time
	Value     float64
}

type SlidingWindowAggregator struct {
	mu          sync.RWMutex
	windowSize  time.Duration
	points      []MetricPoint
	percentiles []float64
}

func NewSlidingWindowAggregator(windowSize time.Duration, percentiles []float64) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		points:      make([]MetricPoint, 0),
		percentiles: percentiles,
	}
}

func (swa *SlidingWindowAggregator) AddPoint(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.points = append(swa.points, MetricPoint{Timestamp: now, Value: value})
	swa.cleanupOldPoints(now)
}

func (swa *SlidingWindowAggregator) cleanupOldPoints(currentTime time.Time) {
	cutoff := currentTime.Add(-swa.windowSize)
	validStart := 0
	for i, point := range swa.points {
		if point.Timestamp.After(cutoff) {
			validStart = i
			break
		}
	}
	swa.points = swa.points[validStart:]
}

func (swa *SlidingWindowAggregator) CalculateStats() (float64, float64, map[float64]float64) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if len(swa.points) == 0 {
		return 0, 0, make(map[float64]float64)
	}

	var sum float64
	values := make([]float64, len(swa.points))
	for i, point := range swa.points {
		sum += point.Value
		values[i] = point.Value
	}
	sort.Float64s(values)

	mean := sum / float64(len(swa.points))

	var varianceSum float64
	for _, v := range values {
		diff := v - mean
		varianceSum += diff * diff
	}
	stdDev := math.Sqrt(varianceSum / float64(len(swa.points)))

	percentileMap := make(map[float64]float64)
	for _, p := range swa.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		index := (p / 100) * float64(len(values)-1)
		lower := int(math.Floor(index))
		upper := int(math.Ceil(index))
		if lower == upper {
			percentileMap[p] = values[lower]
		} else {
			weight := index - float64(lower)
			percentileMap[p] = values[lower]*(1-weight) + values[upper]*weight
		}
	}

	return mean, stdDev, percentileMap
}

func main() {
	aggregator := NewSlidingWindowAggregator(5*time.Minute, []float64{50, 90, 95, 99})

	for i := 0; i < 100; i++ {
		value := float64(i) + 0.5*float64(i%3)
		aggregator.AddPoint(value)
		time.Sleep(100 * time.Millisecond)
	}

	mean, stdDev, percentiles := aggregator.CalculateStats()
	fmt.Printf("Mean: %.2f\n", mean)
	fmt.Printf("Standard Deviation: %.2f\n", stdDev)
	for p, v := range percentiles {
		fmt.Printf("P%.0f: %.2f\n", p, v)
	}
}