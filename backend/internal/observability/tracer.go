package observability

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("your-service-name")

// TracerStart is a helper function to start a new span.
// It's part of the core observability logic.
func TracerStart(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, name)
	span.SetAttributes(attrs...)
	return ctx, span
}

func InitTracer() func() {
	ctx := context.Background()
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("my-service"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		_ = tracerProvider.Shutdown(ctx)
	}
}
