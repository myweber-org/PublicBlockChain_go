
package main

import (
    "log"
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

    httpRequestTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(httpRequestTotal)
}

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

func mainHandler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(50 * time.Millisecond)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func main() {
    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.Handler())
    mux.Handle("/", metricsMiddleware(http.HandlerFunc(mainHandler)))

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal("Server failed:", err)
    }
}package metrics

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	windowSize  time.Duration
	bucketSize  time.Duration
	maxBuckets  int
	buckets     []*Bucket
	currentIdx  int
	startTime   time.Time
	mu          sync.RWMutex
}

type Bucket struct {
	Count int64
	Sum   float64
	Min   float64
	Max   float64
}

func NewSlidingWindow(windowSize, bucketSize time.Duration) *SlidingWindow {
	maxBuckets := int(windowSize/bucketSize) + 1
	sw := &SlidingWindow{
		windowSize: windowSize,
		bucketSize: bucketSize,
		maxBuckets: maxBuckets,
		buckets:    make([]*Bucket, maxBuckets),
		startTime:  time.Now(),
	}
	for i := range sw.buckets {
		sw.buckets[i] = &Bucket{Min: 1<<63 - 1}
	}
	return sw
}

func (sw *SlidingWindow) Add(value float64) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.rotateBuckets()
	bucket := sw.buckets[sw.currentIdx]
	bucket.Count++
	bucket.Sum += value
	if value < bucket.Min {
		bucket.Min = value
	}
	if value > bucket.Max {
		bucket.Max = value
	}
}

func (sw *SlidingWindow) rotateBuckets() {
	now := time.Now()
	elapsed := now.Sub(sw.startTime)
	bucketsPassed := int(elapsed / sw.bucketSize)

	if bucketsPassed > 0 {
		steps := bucketsPassed % sw.maxBuckets
		for i := 0; i < steps; i++ {
			sw.currentIdx = (sw.currentIdx + 1) % sw.maxBuckets
			sw.buckets[sw.currentIdx] = &Bucket{Min: 1<<63 - 1}
		}
		sw.startTime = sw.startTime.Add(time.Duration(bucketsPassed) * sw.bucketSize)
	}
}

func (sw *SlidingWindow) Stats() (count int64, sum, avg, min, max float64) {
	sw.mu.RLock()
	defer sw.mu.RUnlock()

	sw.rotateBuckets()
	var totalCount int64
	var totalSum float64
	currentMin := 1<<63 - 1.0
	currentMax := -1 << 63

	for i, bucket := range sw.buckets {
		if i == sw.currentIdx || bucket.Count > 0 {
			totalCount += bucket.Count
			totalSum += bucket.Sum
			if bucket.Min < currentMin {
				currentMin = bucket.Min
			}
			if bucket.Max > currentMax {
				currentMax = bucket.Max
			}
		}
	}

	if totalCount == 0 {
		return 0, 0, 0, 0, 0
	}
	return totalCount, totalSum, totalSum / float64(totalCount), currentMin, currentMax
}