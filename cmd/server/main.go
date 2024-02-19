package main

import (
	"context"
	"log"

	"github.com/izzanzahrial/skeleton/config"
	db "github.com/izzanzahrial/skeleton/db/sqlc"
	"github.com/izzanzahrial/skeleton/internal/domain/authentication/cache"
	authhandler "github.com/izzanzahrial/skeleton/internal/interface/http/authentication"
	"github.com/izzanzahrial/skeleton/internal/interface/http/handlers"
	"github.com/izzanzahrial/skeleton/internal/interface/http/router"
	userhandler "github.com/izzanzahrial/skeleton/internal/interface/http/user"
	"github.com/izzanzahrial/skeleton/internal/service/authentication"
	"github.com/izzanzahrial/skeleton/internal/service/user"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		panic("failed to load environment variables")
	}

	cfg, err := config.New()
	if err != nil {
		log.Panicf("failed to create config: %v", err)
	}

	conn, err := pgx.Connect(ctx, cfg.Database.URL())
	if err != nil {
		log.Println(cfg.Database.URL())
		log.Panicf("failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	opt, err := redis.ParseURL(cfg.Cache.URL())
	if err != nil {
		log.Panicf("failed to parse URL cache: %v", err)
	}
	rdb := redis.NewClient(opt)

	db := db.New(conn)
	cache := cache.New(rdb)

	authService := authentication.NewService(db, cache)
	authHandler := authhandler.NewHandler(authService)

	userService := user.NewService(db)
	userHandler := userhandler.NewHandler(userService)

	handlers := handlers.NewHandlers(authHandler, userHandler)

	server := echo.New()
	server.Use(middleware.Logger(), middleware.Logger())
	router.MapRoutes(server, handlers)
	if err := server.Start(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}
