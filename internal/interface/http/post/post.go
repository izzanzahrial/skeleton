package post

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/izzanzahrial/skeleton/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type postService interface {
	CreatePost(ctx context.Context, userID int64, title, content string) (model.Post, error)
	GetPostByUserID(ctx context.Context, userID int64) ([]model.Post, error)
	GetPostsFullText(ctx context.Context, limit, offset int, title, content string) ([]model.Post, error)
}

type Handler struct {
	service postService
}

func NewHandler(service postService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreatePost(c echo.Context) error {
	var request CreatPostReq
	if err := c.Bind(&request); err != nil {
		log.Fatalf("failed to bind request: %v", err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		log.Fatalf("failed to validate request: %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	post, err := h.service.CreatePost(context.Background(), request.UserID, request.Title, request.Content)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, post)
}

func (h *Handler) GetPostByUserID(c echo.Context) error {
	var request GetPostByUserIDReq
	if err := c.Bind(&request); err != nil {
		log.Fatalf("failed to bind request: %v", err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		log.Fatalf("failed to validate request: %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	posts, err := h.service.GetPostByUserID(context.Background(), request.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusFound, posts)
}

func (h *Handler) GetPostsFullText(c echo.Context) error {
	var request GetPostsFullTextReq
	if err := c.Bind(&request); err != nil {
		log.Fatalf("failed to bind request: %v", err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		log.Fatalf("failed to validate request: %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	posts, err := h.service.GetPostsFullText(context.Background(), request.Limit, request.Offset, request.Title, request.Content)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusFound, posts)
}
