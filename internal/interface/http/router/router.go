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
	// using native authentication
	e.POST("/login", h.Auth.Login)

	// using google oauth2 authentication
	// temporary name
	e.GET("/google", h.Auth.LoginGoogleOAuth)
	e.GET("/callback", h.Auth.Callback)
	e.GET("/refresh", h.Auth.RefreshToken)
}

func mapUserRoutes(e *echo.Group, h *handlers.Handlers) {
	e.POST("/signup", h.User.Signup)
	e.POST("/signup-admin", h.User.SignUpAdmin, middleware.IsAuthenticated(), middleware.IsAuthorize)
	e.GET("/users/:role", h.User.GetUsersByRole, middleware.IsAuthenticated(), middleware.IsAuthorize)
	e.GET("/users", h.User.GetUsersLikeUsername, middleware.IsAuthenticated(), middleware.IsAuthorize)
	e.DELETE("/users", h.User.DeleteUser)
}
