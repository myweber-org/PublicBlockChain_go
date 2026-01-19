package main

import (
	"log"
	"net/http"
	"time"
)

type MetricsCollector struct {
	latencies []time.Duration
	errors    int
	requests  int
}

func (m *MetricsCollector) RecordLatency(d time.Duration) {
	m.latencies = append(m.latencies, d)
}

func (m *MetricsCollector) RecordError() {
	m.errors++
}

func (m *MetricsCollector) RecordRequest() {
	m.requests++
}

func (m *MetricsCollector) AverageLatency() time.Duration {
	if len(m.latencies) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range m.latencies {
		total += d
	}
	return total / time.Duration(len(m.latencies))
}

func (m *MetricsCollector) ErrorRate() float64 {
	if m.requests == 0 {
		return 0.0
	}
	return float64(m.errors) / float64(m.requests)
}

func main() {
	collector := &MetricsCollector{}

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		collector.RecordRequest()

		// Simulate some processing
		time.Sleep(50 * time.Millisecond)

		// Simulate occasional errors
		if time.Now().UnixNano()%10 == 0 {
			collector.RecordError()
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
		collector.RecordLatency(time.Since(start))
	})

	go func() {
		for {
			time.Sleep(30 * time.Second)
			log.Printf("Metrics - Avg Latency: %v, Error Rate: %.2f%%, Total Requests: %d",
				collector.AverageLatency(),
				collector.ErrorRate()*100,
				collector.requests)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}