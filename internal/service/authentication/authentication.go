package authentication

import (
	"context"
	"errors"
	"log"

	db "github.com/izzanzahrial/skeleton/db/sqlc"
	"github.com/izzanzahrial/skeleton/internal/model"
	pass "github.com/izzanzahrial/skeleton/pkg/password"
	"github.com/izzanzahrial/skeleton/pkg/token"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type authRepo interface {
	GetuserByEmailOrUsername(ctx context.Context, param db.GetuserByEmailOrUsernameParams) (db.User, error)
	CreateUserGoogle(ctx context.Context, param db.CreateUserGoogleParams) (db.User, error)
	GetuserByEmail(ctx context.Context, email string) (db.User, error)
}

type authCache interface {
	SetAuthToken(ctx context.Context, token token.Token) error
}
type Service struct {
	repo  authRepo
	cache authCache
}

type ServiceConfig func(s *Service) error

func NewService(repo authRepo, cache authCache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

// Optional factory pattern
// type ServiceConfig func(s *Service) error

// func NewService(repo authRepo, cfgs ...ServiceConfig) (*Service, error) {
// 	s := &Service{
// 		repo: repo,
// 	}

// 	for _, cfg := range cfgs {
// 		if err := cfg(s); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return s, nil
// }

// func WithRedisCache(c authCache) ServiceConfig {
// 	return func(s *Service) error {
// 		s.cache = c
// 		return nil
// 	}
// }

func (s *Service) GetuserByEmailOrUsername(ctx context.Context, email, username, password string) (model.User, error) {
	param := db.GetuserByEmailOrUsernameParams{
		Email:    email,
		Username: pgtype.Text{String: username, Valid: true},
	}

	user, err := s.repo.GetuserByEmailOrUsername(ctx, param)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Fatalf("Error getting user: %v", err)
		}
		return model.User{}, err
	}

	ok, err := pass.Check(password, user.PasswordHash)
	if !ok || err != nil {
		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Fatalf("Error checking password: %v", err)
		}
		return model.User{}, err
	}

	return model.DBUserToModelUser(user)[0], nil
}

func (s *Service) CreateOrCheckGoogleUser(ctx context.Context, user model.User) (model.User, error) {
	dbUser, err := s.repo.GetuserByEmail(ctx, user.Email)
	if err == nil {
		return model.DBUserToModelUser(dbUser)[0], nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		log.Fatalf("error creating user: %v", err)
		return model.User{}, err
	}

	param := db.CreateUserGoogleParams{
		Email:        user.Email,
		FirstName:    pgtype.Text{String: user.FirstName, Valid: true},
		LastName:     pgtype.Text{String: user.LastName, Valid: true},
		PictureUrl:   pgtype.Text{String: user.PictureUrl, Valid: true},
		RefreshToken: pgtype.Text{String: user.RefreshToken, Valid: true},
		Role:         db.RolesUser,
		Origin:       db.OriginsGoogle,
	}

	dbUser, err = s.repo.CreateUserGoogle(ctx, param)
	if err != nil {
		return model.User{}, err
	}

	return model.DBUserToModelUser(dbUser)[0], nil
}
