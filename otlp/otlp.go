package otlp

import (
	"context"

	internalmetric "github.com/izzanzahrial/skeleton/otlp/metric"
	internaltrace "github.com/izzanzahrial/skeleton/otlp/trace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewMeterAndTraceProvider(ctx context.Context, endpoint string) (*metric.MeterProvider, *trace.TracerProvider, error) {
	traceExp, err := internaltrace.NewGRPCOTLPTraceExporter(context.Background(), endpoint)
	if err != nil {
		return nil, nil, err
	}

	metricExp, err := internalmetric.NewGRPCOTLPMetricExporter(context.Background(), endpoint)
	if err != nil {
		return nil, nil, err
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp, err := internaltrace.NewTraceProvider(traceExp)
	if err != nil {
		return nil, nil, err
	}

	mp, err := internalmetric.NewMeterProvider(metricExp)
	if err != nil {
		return nil, nil, err
	}

	return mp, tp, nil
}
