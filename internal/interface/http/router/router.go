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
	mapPostRoute(v1, h)
}

func mapAuthenticationRoutes(e *echo.Group, h *handlers.Handlers) {
	// using native authentication
	e.POST("/login", h.Auth.Login)

	// using google oauth2 authentication
	e.GET("/google", h.Auth.LoginGoogleOAuth)
	e.GET("/callback", h.Auth.Callback)
	e.GET("/refresh", h.Auth.RefreshToken)

	// using auth0 authentication
	e.GET("/auth0", h.Auth.LoginAuth0)
	e.GET("/callback/auth0", h.Auth.CallbackAuth0)

}

func mapUserRoutes(e *echo.Group, h *handlers.Handlers) {
	e.POST("/signup", h.User.Signup)
	e.POST("/signup-admin", h.User.SignUpAdmin, middleware.IsAuthenticated(), middleware.IsAuthorize)
	e.GET("/users/:role", h.User.GetUsersByRole, middleware.IsAuthenticated(), middleware.IsAuthorize)
	e.GET("/users", h.User.GetUsersLikeUsername, middleware.IsAuthenticated(), middleware.IsAuthorize)
	e.PATCH("/users/:id", h.User.UpdateUser)
	e.DELETE("/users", h.User.DeleteUser, middleware.IsAuthenticated(), middleware.IsAuthorize)
}

func mapPostRoute(e *echo.Group, h *handlers.Handlers) {
	e.POST("/post", h.Post.CreatePost)
	e.GET("/post/:id", h.Post.GetPostByUserID)
	e.GET("/post", h.Post.GetPostsFullText)
}
