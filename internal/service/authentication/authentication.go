package authentication

import (
	"context"

	db "github.com/izzanzahrial/skeleton/db/sqlc"
	"github.com/izzanzahrial/skeleton/internal/model"
	pass "github.com/izzanzahrial/skeleton/pkg/password"
	"github.com/izzanzahrial/skeleton/pkg/token"
)

type authRepo interface {
	GetuserByEmailOrUsername(ctx context.Context, param db.GetuserByEmailOrUsernameParams) (db.User, error)
}

type authCache interface {
	SetAuthToken(ctx context.Context, token token.Token) error
}
type Service struct {
	repo  authRepo
	cache authCache
}

func NewService(repo authRepo, cache authCache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) GetuserByEmailOrUsername(ctx context.Context, email, username, password string) (model.User, string, error) {
	param := db.GetuserByEmailOrUsernameParams{
		Email:    email,
		Username: username,
	}

	user, err := s.repo.GetuserByEmailOrUsername(ctx, param)
	if err != nil {
		return model.User{}, "", err
	}

	ok, err := pass.Check(password, user.PasswordHash)
	if !ok || err != nil {
		return model.User{}, "", err
	}

	// Using token
	// expiry := 24 * time.Hour
	// token, err := token.New(user.ID, expiry)
	// if err != nil {
	// 	return model.User{}, err
	// }

	// if err := s.cache.SetAuthToken(ctx, *token); err != nil {
	// 	return model.User{}, err
	// }

	// Using jwt
	token, err := token.NewJWT(user.ID, model.Roles(user.Role))
	if err != nil {
		return model.User{}, "", err
	}

	return model.DBUserToModelUser(user)[0], token, nil
}
