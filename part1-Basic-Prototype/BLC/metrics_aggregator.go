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
}package metrics

import (
	"sync"
	"time"
)

type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

type Metric struct {
	Name  string
	Type  MetricType
	Value float64
	Tags  map[string]string
	Time  time.Time
}

type SlidingWindow struct {
	windowSize time.Duration
	metrics    []Metric
	mu         sync.RWMutex
}

func NewSlidingWindow(windowSize time.Duration) *SlidingWindow {
	return &SlidingWindow{
		windowSize: windowSize,
		metrics:    make([]Metric, 0),
	}
}

func (sw *SlidingWindow) AddMetric(metric Metric) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	
	metric.Time = time.Now()
	sw.metrics = append(sw.metrics, metric)
	sw.cleanup()
}

func (sw *SlidingWindow) cleanup() {
	cutoff := time.Now().Add(-sw.windowSize)
	validStart := 0
	
	for i, metric := range sw.metrics {
		if metric.Time.After(cutoff) {
			validStart = i
			break
		}
	}
	
	if validStart > 0 {
		sw.metrics = sw.metrics[validStart:]
	}
}

func (sw *SlidingWindow) Aggregate(metricName string, operation string) float64 {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	
	sw.cleanup()
	var result float64
	
	for _, metric := range sw.metrics {
		if metric.Name != metricName {
			continue
		}
		
		switch operation {
		case "sum":
			result += metric.Value
		case "avg":
			result += metric.Value
		case "max":
			if metric.Value > result {
				result = metric.Value
			}
		case "min":
			if len(sw.metrics) == 0 || metric.Value < result {
				result = metric.Value
			}
		}
	}
	
	if operation == "avg" && len(sw.metrics) > 0 {
		count := 0
		for _, metric := range sw.metrics {
			if metric.Name == metricName {
				count++
			}
		}
		if count > 0 {
			result = result / float64(count)
		}
	}
	
	return result
}

func (sw *SlidingWindow) GetMetrics() []Metric {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	
	sw.cleanup()
	return append([]Metric{}, sw.metrics...)
}

type Aggregator struct {
	windows map[string]*SlidingWindow
	mu      sync.RWMutex
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		windows: make(map[string]*SlidingWindow),
	}
}

func (a *Aggregator) RegisterWindow(name string, windowSize time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	a.windows[name] = NewSlidingWindow(windowSize)
}

func (a *Aggregator) AddToWindow(windowName string, metric Metric) {
	a.mu.RLock()
	window, exists := a.windows[windowName]
	a.mu.RUnlock()
	
	if exists {
		window.AddMetric(metric)
	}
}

func (a *Aggregator) GetWindowAggregation(windowName, metricName, operation string) float64 {
	a.mu.RLock()
	window, exists := a.windows[windowName]
	a.mu.RUnlock()
	
	if !exists {
		return 0
	}
	
	return window.Aggregate(metricName, operation)
}