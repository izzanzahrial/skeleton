package trace

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
)

// TracerProvider is an OpenTelemetry TracerProvider.
// It provides Tracers to instrumentation so it can trace operational flow through a system.
func NewTraceProvider(exp trace.SpanExporter) (*trace.TracerProvider, error) {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("skeleton"),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(
		// the exporter that we use to send trace or metrics data to the collector
		trace.WithBatcher(exp),
		trace.WithResource(r),
		// handle rate sampling, by 0.5 means only half of the time trace will be sent
		trace.WithSampler(trace.TraceIDRatioBased(0.5)),
	), nil
}

// List of supported exporters
// https://opentelemetry.io/docs/instrumentation/go/exporters/
// Console Exporter, only for testing
func NewConsoleExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

// OTLP HTTP Exporter
func NewHTTPOTLPExporter(ctx context.Context, endpoint string) (trace.SpanExporter, error) {
	// Change default HTTPS -> HTTP
	insecureOpt := otlptracehttp.WithInsecure()

	endpointOpt := otlptracehttp.WithEndpoint(endpoint)

	return otlptracehttp.New(ctx, insecureOpt, endpointOpt)
}

// OTLP GRPC Exporter
func NewGRPCOTLPTraceExporter(ctx context.Context, endpoint string) (trace.SpanExporter, error) {
	insecureOpt := otlptracegrpc.WithInsecure()

	endpointOpt := otlptracegrpc.WithEndpoint(endpoint)

	return otlptracegrpc.New(ctx, insecureOpt, endpointOpt, otlptracegrpc.WithDialOption(grpc.WithBlock()))
}
