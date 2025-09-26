package middleware

import (
	"context"
	"net/http"
	"time"

	"cyber-go/internal/observability"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation" // You need to import this for trace.Span
	"go.uber.org/zap"
)

// Define a new, unexported type for the context key to avoid collisions.
type contextKey struct{}

var requestIDKey = contextKey{}

// RequestIDKey is a public-facing variable to allow other packages
// to safely retrieve the request ID from the context.
var RequestIDKey = requestIDKey

func ObservabilityMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	tracer := otel.Tracer("cyber-go")
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			ctx, span := tracer.Start(ctx, r.URL.Path)
			defer span.End()

			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.New().String()
			}

			// Use the package-level requestIDKey
			ctx = context.WithValue(ctx, requestIDKey, reqID)
			r = r.WithContext(ctx)

			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			logger.Info("Incoming request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("requestID", reqID),
			)

			next.ServeHTTP(ww, r)

			duration := time.Since(start).Seconds()
			observability.ObserveHTTPRequestDuration(r.URL.Path, r.Method, http.StatusText(ww.statusCode), duration)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// MiddlewareScraper runs in the background and periodically scrapes /metrics
func MiddlewareScraper(interval time.Duration) {
	go func() {
		time.Sleep(3 * time.Second) // wait for server startup
		for {
			observability.ScrapeMetrics()
			time.Sleep(interval)
		}
	}()
}
