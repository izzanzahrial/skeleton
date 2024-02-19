package router

import (
	"github.com/izzanzahrial/skeleton/internal/interface/http/handlers"
	"github.com/izzanzahrial/skeleton/internal/interface/http/middleware"
	"github.com/labstack/echo/v4"
)

func MapRoutes(e *echo.Echo, h *handlers.Handlers) {
	api := e.Group("/api")
	v1 := api.Group("/v1")

	mapAuthenticationRoutes(v1, h)
	mapUserRoutes(v1, h)
}

func mapAuthenticationRoutes(e *echo.Group, h *handlers.Handlers) {
	e.POST("/login", h.Auth.Login)
}

func mapUserRoutes(e *echo.Group, h *handlers.Handlers) {
	e.POST("/signup", h.User.Signup)
	e.POST("/signup-admin", h.User.SignUpAdmin, middleware.IsAuthorize)
	e.GET("/users/:role", h.User.GetUsersByRole, middleware.IsAuthorize)
	e.GET("/users", h.User.GetUsersLikeUsername, middleware.IsAuthorize)
	e.DELETE("/users", h.User.DeleteUser)
}
