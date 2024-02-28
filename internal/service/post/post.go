package post

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	db "github.com/izzanzahrial/skeleton/db/sqlc"
	"github.com/izzanzahrial/skeleton/internal/domain/post/broker"
	"github.com/izzanzahrial/skeleton/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type postRepo interface {
	CreatePost(ctx context.Context, arg db.CreatePostParams) (db.Post, error)
	GetPostByUserID(ctx context.Context, userID int64) ([]db.Post, error)
	GetPostsFullText(ctx context.Context, arg db.GetPostsFullTextParams) ([]db.Post, error)
}

type Service struct {
	repo     postRepo
	producer *broker.Producer
	slog     *slog.Logger
}

func NewService(repo postRepo, producer *broker.Producer, slog *slog.Logger) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
		slog:     slog,
	}
}

func (s *Service) CreatePost(ctx context.Context, userID int64, title, content string) (model.Post, error) {
	user := db.CreatePostParams{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	post, err := s.repo.CreatePost(ctx, user)
	if err != nil {
		s.slog.Error("failed to create post", slog.String("error", err.Error()))
		return model.Post{}, err
	}

	modelPost := model.DBPostToModelPost(post)[0]

	msgPost, err := json.Marshal(modelPost)
	if err != nil {
		s.slog.Error("failed to marshal post", slog.String("error", err.Error()))
		return model.Post{}, err
	}

	if err := s.producer.Publish(ctx, "posts", msgPost); err != nil {
		s.slog.Error("failed to publish post", slog.String("error", err.Error()))
		return modelPost, err
	}

	return modelPost, nil
}

func (s *Service) GetPostByUserID(ctx context.Context, userID int64) ([]model.Post, error) {
	posts, err := s.repo.GetPostByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.slog.Error("failed to get post by used id", slog.String("error", err.Error()))
		}
		return nil, err
	}

	modelPost := model.DBPostToModelPost(posts...)
	return modelPost, nil
}

func (s *Service) GetPostsFullText(ctx context.Context, limit, offset int, keyword string) ([]model.Post, error) {
	var newLimit pgtype.Int4
	if limit <= 0 {
		newLimit.Valid = false
	} else {
		newLimit.Int32 = int32(limit)
		newLimit.Valid = true
	}

	param := db.GetPostsFullTextParams{
		Offset:     int32(offset),
		Keyword:    keyword,
		LimitParam: newLimit,
	}

	posts, err := s.repo.GetPostsFullText(ctx, param)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.slog.Error("failed to get post with keyword", slog.String("error", err.Error()), slog.String("keyword", keyword))
		}
		return nil, err
	}

	modelPost := model.DBPostToModelPost(posts...)
	return modelPost, nil
}
