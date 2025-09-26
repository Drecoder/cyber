// cyber-go/internal/observability/observability.go
package observability

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	// Singleton metrics
	// This is now a HistogramVec to match the labels you want to use.
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	dbQueryDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "db_query_duration_seconds",
		Help: "Duration of DB queries",
	})

	once sync.Once
)

// RegisterMetrics registers metrics only once and logs using Zap
func RegisterMetrics(logger *zap.Logger) {
	once.Do(func() {
		prometheus.MustRegister(httpRequestDuration, dbQueryDuration)
		logger.Info("Metrics successfully registered")
	})
}

// ... ExposeMetricsHandler remains the same ...

// ObserveHTTPRequestDuration now accepts labels for the HistogramVec
func ObserveHTTPRequestDuration(path, method, status string, seconds float64) {
	httpRequestDuration.WithLabelValues(path, method, status).Observe(seconds)
}

// Record DB query duration
func ObserveDBQueryDuration(seconds float64) {
	dbQueryDuration.Observe(seconds)
}
