package authentication

import (
	"context"
	"net/http"

	"github.com/izzanzahrial/skeleton/internal/model"
	"github.com/labstack/echo/v4"
)

type authService interface {
	GetuserByEmailOrUsername(ctx context.Context, email, username, password string) (model.User, string, error)
}

type Handler struct {
	service authService
}

func NewHandler(service authService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, token, err := h.service.GetuserByEmailOrUsername(ctx, email, username, password)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusFound, echo.Map{"user": user, "token": token})
}
