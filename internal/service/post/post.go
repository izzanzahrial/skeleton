package post

import (
	"context"
	"errors"
	"log"

	db "github.com/izzanzahrial/skeleton/db/sqlc"
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
	repo postRepo
}

func NewService(repo postRepo) *Service {
	return &Service{
		repo: repo,
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
		log.Fatalf("failed to create post: %v", err)
		return model.Post{}, err
	}

	modelPost := model.DBPostToModelPost(post)[0]
	return modelPost, nil
}

func (s *Service) GetPostByUserID(ctx context.Context, userID int64) ([]model.Post, error) {
	posts, err := s.repo.GetPostByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Fatalf("failed to get post by user id: %v", err)
		}
		return nil, err
	}

	modelPost := model.DBPostToModelPost(posts...)
	return modelPost, nil
}

func (s *Service) GetPostsFullText(ctx context.Context, limit, offset int, title, content string) ([]model.Post, error) {
	var newLimit pgtype.Int4
	if limit <= 0 {
		newLimit.Valid = false
	} else {
		newLimit.Int32 = int32(limit)
		newLimit.Valid = true
	}

	param := db.GetPostsFullTextParams{
		Offset:     int32(offset),
		Title:      title,
		Content:    content,
		LimitParam: newLimit,
	}

	posts, err := s.repo.GetPostsFullText(ctx, param)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Fatalf("failed to get post by title: %v and content: %v error: %v", title, content, err)
		}
		return nil, err
	}

	modelPost := model.DBPostToModelPost(posts...)
	return modelPost, nil
}
