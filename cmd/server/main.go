package main

import (
	"context"
	"log"
	"os"

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
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic("failed to load environment variables")
	}

	dbCfg, err := config.NewDatabase()
	if err != nil {
		log.Panic("failed to initialize database configuration")
	}

	conn, err := pgx.Connect(context.Background(), dbCfg.URL())
	if err != nil {
		log.Println(dbCfg.URL())
		log.Panicf("failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	redisCfg, err := config.NewCache()
	if err != nil {
		log.Panic("failed to initialize cache configuration")
	}

	opt, err := redis.ParseURL(redisCfg.URL())
	if err != nil {
		log.Println(redisCfg.URL())
		log.Panicf("failed to parse URL cache: %v", err)
	}
	rdb := redis.NewClient(opt)

	db := db.New(conn)
	cache := cache.New(rdb)

	auht0, err := auth0.New()
	if err != nil {
		log.Panicf("failed to create auth0: %v", err)
	}

	producer, err := broker.NewProducer()
	if err != nil {
		log.Panicf("failed to create producer: %v", err)
	}

	authService := authentication.NewService(db, cache)
	authHandler := authhandler.NewHandler(authService, auht0)

	userService := user.NewService(db)
	userHandler := userhandler.NewHandler(userService)

	postService := post.NewService(db, producer)
	postHandler := posthandler.NewHandler(postService)

	handlers := handlers.NewHandlers(authHandler, userHandler, postHandler)

	cv, err := pkgvalidator.New()
	if err != nil {
		log.Panic(err)
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
		log.Panicf("failed to start server: %s", err.Error())
	}
}
