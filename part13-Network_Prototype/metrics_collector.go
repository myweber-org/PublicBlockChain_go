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
		
		log.Printf("Request processed: %s %s - Status: %d, Duration: %v", 
			r.Method, r.URL.Path, recorder.statusCode, duration)
	})
}

func (mc *MetricsCollector) GetMetrics() (int, int, time.Duration) {
	return mc.requestCount, mc.errorCount, mc.totalLatency
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
		reqCount, errCount, totalLatency := collector.GetMetrics()
		avgLatency := time.Duration(0)
		if reqCount > 0 {
			avgLatency = totalLatency / time.Duration(reqCount)
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"request_count": ` + string(rune(reqCount)) + `,
			"error_count": ` + string(rune(errCount)) + `,
			"average_latency_ns": ` + string(rune(avgLatency.Nanoseconds())) + `
		}`))
	})
	
	handler := collector.Middleware(mux)
	
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}