package cache

import (
	"context"
	"fmt"

	"github.com/izzanzahrial/skeleton/pkg/token"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	rdb *redis.Client
}

func New(redis *redis.Client) *Repository {
	return &Repository{rdb: redis}
}

func (r *Repository) SetAuthToken(ctx context.Context, token token.Token) error {
	if err := r.rdb.Set(ctx, string(token.Hash), token.User.ID, token.Expiry).Err(); err != nil {
		return fmt.Errorf("failed to set auth token into redis cache: %v", err)
	}

	return nil
}
