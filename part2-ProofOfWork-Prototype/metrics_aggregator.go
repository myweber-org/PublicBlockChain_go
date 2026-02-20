
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
	metric := Metric{Timestamp: now, Value: value}
	swa.metrics = append(swa.metrics, metric)

	swa.cleanupOldMetrics(now)
	if len(swa.metrics) > swa.maxSamples {
		swa.metrics = swa.metrics[len(swa.metrics)-swa.maxSamples:]
	}
}

func (swa *SlidingWindowAggregator) cleanupOldMetrics(currentTime time.Time) {
	cutoff := currentTime.Add(-swa.windowSize)
	i := 0
	for i < len(swa.metrics) && swa.metrics[i].Timestamp.Before(cutoff) {
		i++
	}
	if i > 0 {
		swa.metrics = swa.metrics[i:]
	}
}

func (swa *SlidingWindowAggregator) GetPercentile(p float64) (float64, error) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if len(swa.metrics) == 0 {
		return 0, fmt.Errorf("no metrics available")
	}

	values := make([]float64, len(swa.metrics))
	for i, m := range swa.metrics {
		values[i] = m.Value
	}
	sort.Float64s(values)

	if p <= 0 {
		return values[0], nil
	}
	if p >= 100 {
		return values[len(values)-1], nil
	}

	index := (p / 100) * float64(len(values)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return values[lower], nil
	}

	weight := index - float64(lower)
	return values[lower]*(1-weight) + values[upper]*weight, nil
}

func (swa *SlidingWindowAggregator) GetStats() (min, max, avg float64, count int) {
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

	avg = sum / float64(len(swa.metrics))
	count = len(swa.metrics)
	return
}

func main() {
	aggregator := NewSlidingWindowAggregator(5*time.Minute, 1000)

	for i := 0; i < 50; i++ {
		value := 10.0 + float64(i%20) + (float64(i) / 10.0)
		aggregator.AddMetric(value)
		time.Sleep(100 * time.Millisecond)
	}

	min, max, avg, count := aggregator.GetStats()
	fmt.Printf("Stats - Min: %.2f, Max: %.2f, Avg: %.2f, Count: %d\n", min, max, avg, count)

	p95, err := aggregator.GetPercentile(95)
	if err != nil {
		fmt.Printf("Error calculating percentile: %v\n", err)
	} else {
		fmt.Printf("95th percentile: %.2f\n", p95)
	}

	p50, err := aggregator.GetPercentile(50)
	if err != nil {
		fmt.Printf("Error calculating percentile: %v\n", err)
	} else {
		fmt.Printf("50th percentile: %.2f\n", p50)
	}
}