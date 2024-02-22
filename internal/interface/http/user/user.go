package user

import (
	"context"
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

	var request SignUpUserReq
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.CreateUser(ctx, request.Email, request.Username, request.Password)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusFound, user)
}

func (h *Handler) SignUpAdmin(c echo.Context) error {
	ctx := c.Request().Context()

	var request SignUpUserReq
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.CreateAdmin(ctx, request.Email, request.Username, request.Password)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusFound, user)
}

func (h *Handler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	var request GetUserReq
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.GetUser(ctx, int64(request.ID))
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

	if err := c.Validate(&param); err != nil {
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

	var request GetUsersLikeUsernameReq
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	users, err := h.service.GetUsersLikeUsername(ctx, request.Username, int32(request.Limit), int32(request.Offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	var request DeleteUserReq
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.service.DeleteUser(ctx, int64(request.ID)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}
