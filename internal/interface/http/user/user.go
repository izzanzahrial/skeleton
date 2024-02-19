package user

import (
	"context"
	"log"
	"net/http"

	"github.com/izzanzahrial/skeleton/internal/model"
	"github.com/labstack/echo/v4"
)

type userService interface {
	CreateUser(ctx context.Context, email, username, password string) (model.User, error)
	CreateAdmin(ctx context.Context, email, username, password string) (model.User, error)
	GetUser(ctx context.Context, id int64) (model.User, error)
	GetUsersByRole(ctx context.Context, role model.Roles, limit, offset int32) ([]model.User, error)
	GetUsersLikeUsername(ctx context.Context, username string, limit, offset int32) ([]model.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type Handler struct {
	service userService
}

func NewHandler(service userService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Signup(c echo.Context) error {
	ctx := c.Request().Context()

	email := c.FormValue("email")
	username := c.FormValue("username")
	log.Println(username)
	password := c.FormValue("password")

	user, err := h.service.CreateUser(ctx, email, username, password)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusFound, user)
}

func (h *Handler) SignUpAdmin(c echo.Context) error {
	ctx := c.Request().Context()

	email := c.FormValue("email")
	username := c.FormValue("username")
	log.Println(username)
	password := c.FormValue("password")

	user, err := h.service.CreateAdmin(ctx, email, username, password)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusFound, user)
}

func (h *Handler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	var param GetUserReq
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.GetUser(ctx, int64(param.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUsersByRole(c echo.Context) error {
	ctx := c.Request().Context()

	var param GetUsersByRoleReq
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	users, err := h.service.GetUsersByRole(ctx, model.Roles(param.Role), int32(param.Limit), int32(param.Offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUsersLikeUsername(c echo.Context) error {
	ctx := c.Request().Context()

	var param GetUsersLikeUsernameReq
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	users, err := h.service.GetUsersLikeUsername(ctx, param.Username, int32(param.Limit), int32(param.Offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	var param DeleteUserReq
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.service.DeleteUser(ctx, int64(param.ID)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}
