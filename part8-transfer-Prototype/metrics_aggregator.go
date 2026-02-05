package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize   time.Duration
	dataPoints   []float64
	timestamps   []time.Time
	mu           sync.RWMutex
	percentiles  []float64
}

func NewAggregator(windowSize time.Duration, percentiles []float64) *Aggregator {
	return &Aggregator{
		windowSize:  windowSize,
		dataPoints:  make([]float64, 0),
		timestamps:  make([]time.Time, 0),
		percentiles: percentiles,
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

func (a *Aggregator) GetStats() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	stats := make(map[string]float64)
	if len(a.dataPoints) == 0 {
		return stats
	}

	stats["count"] = float64(len(a.dataPoints))
	stats["min"], stats["max"], stats["sum"] = a.calculateBasicStats()

	sortedPoints := make([]float64, len(a.dataPoints))
	copy(sortedPoints, a.dataPoints)
	sort.Float64s(sortedPoints)

	for _, p := range a.percentiles {
		key := formatPercentileKey(p)
		stats[key] = calculatePercentile(sortedPoints, p)
	}

	return stats
}

func (a *Aggregator) calculateBasicStats() (min, max, sum float64) {
	min = a.dataPoints[0]
	max = a.dataPoints[0]
	sum = 0

	for _, v := range a.dataPoints {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	return
}

func calculatePercentile(sortedData []float64, percentile float64) float64 {
	if len(sortedData) == 0 {
		return 0
	}

	index := (percentile / 100) * float64(len(sortedData)-1)
	lower := int(index)
	upper := lower + 1

	if upper >= len(sortedData) {
		return sortedData[lower]
	}

	weight := index - float64(lower)
	return sortedData[lower]*(1-weight) + sortedData[upper]*weight
}

func formatPercentileKey(p float64) string {
	return "p" + strings.Replace(fmt.Sprintf("%.1f", p), ".", "_", -1)
}