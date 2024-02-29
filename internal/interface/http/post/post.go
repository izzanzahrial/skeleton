package post

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/izzanzahrial/skeleton/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/izzanzahrial/skeleton/internal/interface/http/post")

type postService interface {
	CreatePost(ctx context.Context, userID int64, title, content string) (model.Post, error)
	GetPostByUserID(ctx context.Context, userID int64) ([]model.Post, error)
	GetPostsFullText(ctx context.Context, limit, offset int, keyword string) ([]model.Post, error)
}

type Handler struct {
	service postService
	slog    *slog.Logger
}

func NewHandler(service postService, slog *slog.Logger) *Handler {
	return &Handler{service: service, slog: slog}
}

func (h *Handler) CreatePost(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span := tracer.Start(ctx, "post.CreatePost")
	defer span.End()

	var request CreatPostReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	post, err := h.service.CreatePost(ctx, request.UserID, request.Title, request.Content)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, post)
}

func (h *Handler) GetPostByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span := tracer.Start(ctx, "post.GetPostByUserID")
	defer span.End()

	var request GetPostByUserIDReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	posts, err := h.service.GetPostByUserID(ctx, request.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusFound, posts)
}

func (h *Handler) GetPostsFullText(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span := tracer.Start(ctx, "post.GetPostsFullText")
	defer span.End()

	var request GetPostsFullTextReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("failed to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	posts, err := h.service.GetPostsFullText(ctx, request.Limit, request.Offset, request.Keyword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusFound, posts)
}
