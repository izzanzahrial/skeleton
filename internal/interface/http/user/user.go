package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/izzanzahrial/skeleton/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type userService interface {
	CreateUser(ctx context.Context, email, username, password string) (model.User, error)
	CreateAdmin(ctx context.Context, email, username, password string) (model.User, error)
	GetUser(ctx context.Context, id int64) (model.User, error)
	GetUsersByRole(ctx context.Context, role model.Roles, limit, offset int32) ([]model.User, error)
	GetUsersLikeUsername(ctx context.Context, username string, limit, offset int32) ([]model.User, error)
	UpdateUser(ctx context.Context, id int64, email, username, password *string) (model.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type Handler struct {
	service userService
	slog    *slog.Logger
}

func NewHandler(service userService, slog *slog.Logger) *Handler {
	return &Handler{service: service, slog: slog}
}

func (h *Handler) Signup(c echo.Context) error {
	var request SignUpUserReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.CreateUser(context.Background(), request.Email, request.Username, request.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) SignUpAdmin(c echo.Context) error {
	var request SignUpUserReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.CreateAdmin(context.Background(), request.Email, request.Username, request.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUser(c echo.Context) error {
	var request GetUserReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.GetUser(context.Background(), int64(request.ID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusFound, user)
}

func (h *Handler) GetUsersByRole(c echo.Context) error {
	var request GetUsersByRoleReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	users, err := h.service.GetUsersByRole(context.Background(), model.Roles(request.Role), int32(request.Limit), int32(request.Offset))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUsersLikeUsername(c echo.Context) error {
	var request GetUsersLikeUsernameReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	users, err := h.service.GetUsersLikeUsername(context.Background(), request.Username, int32(request.Limit), int32(request.Offset))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	var request DeleteUserReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.service.DeleteUser(context.Background(), int64(request.ID)); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	var request UpdateUserReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.UpdateUser(context.Background(), int64(request.ID), request.Email, request.Username, request.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, user)
}
