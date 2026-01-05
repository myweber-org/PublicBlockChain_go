package aggregator

import (
	"sync"
	"time"
)

type Metric struct {
	Timestamp time.Time
	Value     float64
}

type SlidingWindowAggregator struct {
	windowSize  time.Duration
	metrics     []Metric
	mu          sync.RWMutex
	subscribers []chan float64
}

func NewSlidingWindowAggregator(windowSize time.Duration) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		metrics:     make([]Metric, 0),
		subscribers: make([]chan float64, 0),
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.metrics = append(swa.metrics, Metric{Timestamp: now, Value: value})
	swa.cleanupOldMetrics(now)
	swa.notifySubscribers()
}

func (swa *SlidingWindowAggregator) cleanupOldMetrics(currentTime time.Time) {
	cutoff := currentTime.Add(-swa.windowSize)
	validStart := 0
	for i, metric := range swa.metrics {
		if metric.Timestamp.After(cutoff) {
			validStart = i
			break
		}
	}
	swa.metrics = swa.metrics[validStart:]
}

func (swa *SlidingWindowAggregator) GetAverage() float64 {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if len(swa.metrics) == 0 {
		return 0.0
	}

	var sum float64
	for _, metric := range swa.metrics {
		sum += metric.Value
	}
	return sum / float64(len(swa.metrics))
}

func (swa *SlidingWindowAggregator) Subscribe() <-chan float64 {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	ch := make(chan float64, 10)
	swa.subscribers = append(swa.subscribers, ch)
	return ch
}

func (swa *SlidingWindowAggregator) notifySubscribers() {
	avg := swa.GetAverage()
	for _, ch := range swa.subscribers {
		select {
		case ch <- avg:
		default:
		}
	}
}

func (swa *SlidingWindowAggregator) Stop() {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	for _, ch := range swa.subscribers {
		close(ch)
	}
	swa.subscribers = nil
}