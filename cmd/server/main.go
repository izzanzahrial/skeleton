package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/izzanzahrial/skeleton/config"
	db "github.com/izzanzahrial/skeleton/db/sqlc"
	"github.com/izzanzahrial/skeleton/internal/domain/authentication/cache"
	"github.com/izzanzahrial/skeleton/internal/domain/post/broker"
	"github.com/izzanzahrial/skeleton/internal/interface/http/auth0"
	authhandler "github.com/izzanzahrial/skeleton/internal/interface/http/authentication"
	"github.com/izzanzahrial/skeleton/internal/interface/http/handlers"
	posthandler "github.com/izzanzahrial/skeleton/internal/interface/http/post"
	"github.com/izzanzahrial/skeleton/internal/interface/http/router"
	userhandler "github.com/izzanzahrial/skeleton/internal/interface/http/user"
	"github.com/izzanzahrial/skeleton/internal/service/authentication"
	"github.com/izzanzahrial/skeleton/internal/service/post"
	"github.com/izzanzahrial/skeleton/internal/service/user"
	pkgvalidator "github.com/izzanzahrial/skeleton/pkg/validator"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	exp, err := newOTLPExporter(context.Background())
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// TODO: move this somewhere else
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to determine working directory")
	}

	replacer := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			if file, ok := strings.CutPrefix(source.File, wd); ok {
				source.File = file
			}
		}
		return a
	}

	options := &slog.HandlerOptions{
		Level:       logLevel(level),
		ReplaceAttr: replacer,
	}
	if level == "debug" {
		options.AddSource = true
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, options))
	slog.SetDefault(logger)

	dbCfg, err := config.NewDatabase()
	if err != nil {
		slog.Warn("failed to initialize database configuration")
	}

	conn, err := pgx.Connect(context.Background(), dbCfg.URL())
	if err != nil {
		slog.Warn("failed to connect to database", slog.String("url", dbCfg.URL()), slog.String("error", err.Error()))
	}
	defer conn.Close(context.Background())

	redisCfg, err := config.NewCache()
	if err != nil {
		slog.Warn("failed to initialize cache configuration")
	}

	opt, err := redis.ParseURL(redisCfg.URL())
	if err != nil {
		slog.Warn("failed to parse URL cache", slog.String("url", redisCfg.URL()), slog.String("error", err.Error()))
	}
	rdb := redis.NewClient(opt)

	db := db.New(conn)
	cache := cache.New(rdb)

	auht0, err := auth0.New()
	if err != nil {
		slog.Warn("failed to create auth0", slog.String("error", err.Error()))
	}

	producer, err := broker.NewProducer()
	if err != nil {
		slog.Warn("failed to create producer", slog.String("error", err.Error()))
	}

	authService := authentication.NewService(db, cache, logger)
	authHandler := authhandler.NewHandler(authService, auht0, logger)

	userService := user.NewService(db, logger)
	userHandler := userhandler.NewHandler(userService, logger)

	postService := post.NewService(db, producer, logger)
	postHandler := posthandler.NewHandler(postService, logger)

	handlers := handlers.NewHandlers(authHandler, userHandler, postHandler)

	cv, err := pkgvalidator.New()
	if err != nil {
		slog.Warn("failed to create validator", slog.String("error", err.Error()))
	}

	server := echo.New()
	server.Use(middleware.Logger())
	server.Validator = cv

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.MapRoutes(server, handlers)
	if err := server.Start(":" + port); err != nil {
		slog.Warn("failed to start server", slog.String("error", err.Error()))
	}
}

func logLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// List of supported exporters
// https://opentelemetry.io/docs/instrumentation/go/exporters/

// Console Exporter, only for testing
func newConsoleExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

// OTLP Exporter
func newOTLPExporter(ctx context.Context) (trace.SpanExporter, error) {
	// Change default HTTPS -> HTTP
	insecureOpt := otlptracehttp.WithInsecure()

	// Update default OTLP reciver endpoint
	// TODO: create oltp endpoint for grafana
	endpointOpt := otlptracehttp.WithEndpoint("")

	return otlptracehttp.New(ctx, insecureOpt, endpointOpt)
}

// TracerProvider is an OpenTelemetry TracerProvider.
// It provides Tracers to instrumentation so it can trace operational flow through a system.
func newTraceProvider(exp trace.SpanExporter) *trace.TracerProvider {
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
		panic(err)
	}

	return trace.NewTracerProvider(
		// the exporter that we use to send trace or metrics data to the collector
		trace.WithBatcher(exp),
		trace.WithResource(r),
		// handle rate sampling, by 0.5 means only half of the time trace will be sent
		trace.WithSampler(trace.TraceIDRatioBased(0.5)))
}
