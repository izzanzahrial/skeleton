package main

import (
	"context"
	"log"

	"github.com/izzanzahrial/skeleton/config"
	db "github.com/izzanzahrial/skeleton/db/sqlc"
	"github.com/izzanzahrial/skeleton/internal/domain/authentication/cache"
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
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Panic("failed to load environment variables")
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
		log.Println(cfg.Cache.URL())
		log.Panicf("failed to parse URL cache: %v", err)
	}
	rdb := redis.NewClient(opt)

	db := db.New(conn)
	cache := cache.New(rdb)

	auht0, err := auth0.New()
	if err != nil {
		log.Panicf("failed to create auth0: %v", err)
	}

	authService := authentication.NewService(db, cache)
	authHandler := authhandler.NewHandler(authService, auht0)

	userService := user.NewService(db)
	userHandler := userhandler.NewHandler(userService)

	postService := post.NewService(db)
	postHandler := posthandler.NewHandler(postService)

	handlers := handlers.NewHandlers(authHandler, userHandler, postHandler)

	cv, err := pkgvalidator.New()
	if err != nil {
		log.Panic(err)
	}

	server := echo.New()
	server.Use(middleware.Logger())
	server.Validator = cv

	router.MapRoutes(server, handlers)
	if err := server.Start(":" + cfg.Port); err != nil {
		log.Panicf("failed to start server: %s", err.Error())
	}
}
