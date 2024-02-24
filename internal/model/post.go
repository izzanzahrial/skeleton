package model

import (
	"time"

	db "github.com/izzanzahrial/skeleton/db/sqlc"
)

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

func DBPostToModelPost(posts ...db.Post) []Post {
	var modelPosts []Post

	for _, p := range posts {
		modelPosts = append(modelPosts, Post{
			ID:        p.ID,
			UserID:    p.UserID,
			CreatedAt: p.CreatedAt.Time,
			UpdatedAt: p.UpdatedAt.Time,
			DeletedAt: p.DeletedAt.Time,
			Title:     p.Title,
			Content:   p.Content,
		})
	}

	return modelPosts
}
