
package main

import (
	"log"
	"net/http"
	"time"
)

type MetricsCollector struct {
	requestCount    int
	errorCount      int
	totalLatency    time.Duration
	statusCodeCount map[int]int
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		statusCodeCount: make(map[int]int),
	}
}

func (mc *MetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		mc.recordRequest(recorder.statusCode, duration)
	})
}

func (mc *MetricsCollector) recordRequest(status int, latency time.Duration) {
	mc.requestCount++
	mc.totalLatency += latency
	mc.statusCodeCount[status]++
	
	if status >= 400 {
		mc.errorCount++
	}
}

func (mc *MetricsCollector) GetAverageLatency() time.Duration {
	if mc.requestCount == 0 {
		return 0
	}
	return mc.totalLatency / time.Duration(mc.requestCount)
}

func (mc *MetricsCollector) GetErrorRate() float64 {
	if mc.requestCount == 0 {
		return 0
	}
	return float64(mc.errorCount) / float64(mc.requestCount)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func main() {
	collector := NewMetricsCollector()
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request processed"))
	})
	
	http.Handle("/", collector.Middleware(handler))
	
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			log.Printf("Requests: %d, Avg Latency: %v, Error Rate: %.2f%%", 
				collector.requestCount, 
				collector.GetAverageLatency(), 
				collector.GetErrorRate()*100)
		}
	}()
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}