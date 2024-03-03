package user

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var tracer = otel.Tracer("github.com/izzanzahrial/skeleton/internal/interface/http/user")
var meter = otel.Meter("github.com/izzanzahrial/skeleton/internal/interface/http/auth")

var signUpCounter, _ = meter.Int64Counter(
	"signup.counter",
	metric.WithDescription("number of API calls to signup handler"),
	metric.WithUnit("{calls}"),
)

var signUpDuration, _ = meter.Float64Histogram(
	"signup.duration",
	metric.WithDescription("the duration of the signup handler"),
	metric.WithUnit("s"),
)

var signUpAdminCounter, _ = meter.Int64Counter(
	"signup.admin.counter",
	metric.WithDescription("number of API calls to signup admin handler"),
	metric.WithUnit("{calls}"),
)

var signUpAdminDuration, _ = meter.Float64Histogram(
	"signup.admin.duration",
	metric.WithDescription("the duration of the signup admin handler"),
	metric.WithUnit("s"),
)

var getUserCounter, _ = meter.Int64Counter(
	"getUser.counter",
	metric.WithDescription("number of API calls to get user handler"),
	metric.WithUnit("{calls}"),
)

var getUserDuration, _ = meter.Float64Histogram(
	"getUser.duration",
	metric.WithDescription("the duration of the get user handler"),
	metric.WithUnit("s"),
)

var getUserByRoleCounter, _ = meter.Int64Counter(
	"getUserByRole.counter",
	metric.WithDescription("number of API calls to get user by role handler"),
	metric.WithUnit("{calls}"),
)

var getUserByRoleDuration, _ = meter.Float64Histogram(
	"getUserByRole.duration",
	metric.WithDescription("the duration of the get user by role handler"),
	metric.WithUnit("s"),
)

var getUsersLikeUsernameCounter, _ = meter.Int64Counter(
	"getUsersLikeUsername.counter",
	metric.WithDescription("number of API calls to get user like username handler"),
	metric.WithUnit("{calls}"),
)

var getUsersLikeUsernameDuration, _ = meter.Float64Histogram(
	"getUsersLikeUsername.duration",
	metric.WithDescription("the duration of the get user like username handler"),
	metric.WithUnit("s"),
)

var deleteUserCounter, _ = meter.Int64Counter(
	"delete.user.counter",
	metric.WithDescription("number of API calls to delete user handler"),
	metric.WithUnit("{calls}"),
)

var deleteUserDuration, _ = meter.Float64Histogram(
	"delete.user.duration",
	metric.WithDescription("the duration of the delete user handler"),
	metric.WithUnit("s"),
)

var updateUserCounter, _ = meter.Int64Counter(
	"update.user.counter",
	metric.WithDescription("number of API calls to update user handler"),
	metric.WithUnit("{calls}"),
)

var updateUserDuration, _ = meter.Float64Histogram(
	"update.user.duration",
	metric.WithDescription("the duration of the update user handler"),
	metric.WithUnit("s"),
)
