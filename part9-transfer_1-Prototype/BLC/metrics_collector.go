
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
}

func (mc *MetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		mc.requestCount++
		mc.totalLatency += duration
		
		if recorder.statusCode >= 400 {
			mc.errorCount++
		}
		
		log.Printf("Request: %s %s - Status: %d - Duration: %v", 
			r.Method, r.URL.Path, recorder.statusCode, duration)
	})
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
	return float64(mc.errorCount) / float64(mc.requestCount) * 100
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
	collector := &MetricsCollector{}
	mux := http.NewServeMux()
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"requests":` + string(rune(collector.requestCount)) + 
			`,"avg_latency_ms":` + string(rune(collector.GetAverageLatency().Milliseconds())) +
			`,"error_rate":` + string(rune(collector.GetErrorRate())) + `}`))
	})
	
	handler := collector.Middleware(mux)
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}