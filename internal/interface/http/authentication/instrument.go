package authentication

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var tracer = otel.Tracer("github.com/izzanzahrial/skeleton/internal/interface/http/auth")
var meter = otel.Meter("github.com/izzanzahrial/skeleton/internal/interface/http/auth")

var loginCounter, _ = meter.Int64Counter(
	"login.counter",
	metric.WithDescription("number of API calls to login handler"),
	metric.WithUnit("{calls}"),
)

var loginDuration, _ = meter.Float64Histogram(
	"login.duration",
	metric.WithDescription("the duration of the login handler"),
	metric.WithUnit("s"),
)
