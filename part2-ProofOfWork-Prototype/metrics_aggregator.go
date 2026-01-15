
package aggregator

import (
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
	dataPoints  []MetricPoint
	mu          sync.RWMutex
	percentiles []float64
}

func NewSlidingWindowAggregator(windowSize time.Duration, percentiles []float64) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		dataPoints:  make([]MetricPoint, 0),
		percentiles: percentiles,
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.dataPoints = append(swa.dataPoints, MetricPoint{
		Value:     value,
		Timestamp: now,
	})
	swa.cleanupOldPoints(now)
}

func (swa *SlidingWindowAggregator) cleanupOldPoints(currentTime time.Time) {
	cutoff := currentTime.Add(-swa.windowSize)
	validStart := 0
	for i, point := range swa.dataPoints {
		if point.Timestamp.After(cutoff) {
			validStart = i
			break
		}
	}
	swa.dataPoints = swa.dataPoints[validStart:]
}

func (swa *SlidingWindowAggregator) GetAggregatedMetrics() map[string]float64 {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if len(swa.dataPoints) == 0 {
		return make(map[string]float64)
	}

	values := make([]float64, len(swa.dataPoints))
	for i, point := range swa.dataPoints {
		values[i] = point.Value
	}

	results := make(map[string]float64)
	results["count"] = float64(len(values))
	results["min"], results["max"], results["avg"] = calculateBasicStats(values)

	if len(values) > 0 {
		sortedValues := make([]float64, len(values))
		copy(sortedValues, values)
		sort.Float64s(sortedValues)

		for _, p := range swa.percentiles {
			if p >= 0 && p <= 100 {
				key := formatPercentileKey(p)
				results[key] = calculatePercentile(sortedValues, p)
			}
		}
	}

	return results
}

func calculateBasicStats(values []float64) (min, max, avg float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}

	min = values[0]
	max = values[0]
	sum := 0.0

	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}

	avg = sum / float64(len(values))
	return min, max, avg
}

func calculatePercentile(sortedValues []float64, percentile float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}

	index := (percentile / 100) * float64(len(sortedValues)-1)
	lowerIndex := int(index)
	upperIndex := lowerIndex + 1

	if upperIndex >= len(sortedValues) {
		return sortedValues[lowerIndex]
	}

	weight := index - float64(lowerIndex)
	return sortedValues[lowerIndex]*(1-weight) + sortedValues[upperIndex]*weight
}

func formatPercentileKey(percentile float64) string {
	return "p" + strconv.FormatFloat(percentile, 'f', -1, 64)
}