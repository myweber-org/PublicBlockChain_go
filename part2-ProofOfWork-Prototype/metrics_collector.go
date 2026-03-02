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
}package main

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path", "status"},
    )

    httpRequestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_request_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(httpRequestCount)
}

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        defer func() {
            duration := time.Since(start).Seconds()
            status := http.StatusText(rw.statusCode)
            
            httpRequestDuration.WithLabelValues(
                r.Method,
                r.URL.Path,
                status,
            ).Observe(duration)
            
            httpRequestCount.WithLabelValues(
                r.Method,
                r.URL.Path,
                status,
            ).Inc()
        }()

        next.ServeHTTP(rw, r)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    mux.Handle("/metrics", promhttp.Handler())
    
    handler := metricsMiddleware(mux)
    
    http.ListenAndServe(":8080", handler)
}
package main

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path", "status"},
    )

    httpRequestTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(rw, r)

        duration := time.Since(start).Seconds()
        status := http.StatusText(rw.statusCode)

        httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
        httpRequestTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    mux.Handle("/metrics", promhttp.Handler())

    handler := metricsMiddleware(mux)
    http.ListenAndServe(":8080", handler)
}