package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/exaring/otelpgx"
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
	"github.com/izzanzahrial/skeleton/otlp"
	pkgvalidator "github.com/izzanzahrial/skeleton/pkg/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	otlpEndpoint := os.Getenv("OTEL_RECEIVER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4317"
	}

	mp, tp, err := otlp.NewMeterAndTraceProvider(context.Background(), otlpEndpoint)
	if err != nil {
		log.Fatalf("failed to create meter and trace provider: %v", err)
	}
	// TODO: move this to graceful shutdown function
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("error shutting down tracer provider: %v", err)
		}
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Printf("error shutting down metric provider: %v", err)
		}
	}()

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
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

	// using pgx
	// conn, err := pgx.Connect(context.Background(), dbCfg.URL())
	// if err != nil {
	// 	slog.Warn("failed to connect to database", slog.String("url", dbCfg.URL()), slog.String("error", err.Error()))
	// }
	// defer conn.Close(context.Background())

	// using pgxpool + opentelemetry auto instrumentation library for pgx : https://github.com/exaring/otelpgx
	cfg, err := pgxpool.ParseConfig(dbCfg.URL())
	if err != nil {
		log.Fatalf("failed to parse database config: %v", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()
	// cfg.MaxConns = 200

	conn, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to create database connection: %v", err)
	}

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
	// add echo instrumentation library https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/labstack/echo
	server.Use(otelecho.Middleware("skeleton-service"), middleware.Logger())
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
