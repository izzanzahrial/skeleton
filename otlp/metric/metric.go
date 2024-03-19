package metric

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func NewMeterProvider(exp metric.Exporter) (*metric.MeterProvider, error) {
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

	return metric.NewMeterProvider(
		metric.WithResource(r),
		metric.WithReader(
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.NewPeriodicReader(exp, metric.WithInterval(3*time.Second)),
		),
	), nil
}

func NewGRPCOTLPMetricExporter(ctx context.Context, endpoint string) (metric.Exporter, error) {
	insecureOpt := otlpmetricgrpc.WithInsecure()

	endpointOpt := otlpmetricgrpc.WithEndpoint(endpoint)

	return otlpmetricgrpc.New(ctx, insecureOpt, endpointOpt)
}
