package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	db "github.com/izzanzahrial/skeleton/db/sqlc"
	"github.com/izzanzahrial/skeleton/internal/model"
	pass "github.com/izzanzahrial/skeleton/pkg/password"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type userRepo interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, id int64) (db.User, error)
	GetUsersByRole(ctx context.Context, arg db.GetUsersByRoleParams) ([]db.User, error)
	GetUsersLikeUsername(ctx context.Context, arg db.GetUsersLikeUsernameParams) ([]db.User, error)
	GetuserByEmailOrUsername(ctx context.Context, arg db.GetuserByEmailOrUsernameParams) (db.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type Service struct {
	repo userRepo
}

func NewService(repo userRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, email, username, password string) (model.User, error) {
	passHash, err := pass.Generate(password)
	if err != nil {
		log.Fatalf("failed to generate password hash: %v", err)
		return model.User{}, err
	}

	user := db.CreateUserParams{
		Email:        email,
		Username:     pgtype.Text{String: username, Valid: true},
		PasswordHash: passHash,
		Role:         db.RolesUser,
		Origin:       db.OriginsNative,
	}

	newUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
		return model.User{}, err
	}

	modelUser := model.DBUserToModelUser(newUser)[0]
	return modelUser, nil
}

func (s *Service) CreateAdmin(ctx context.Context, email, username, password string) (model.User, error) {
	passHash, err := pass.Generate(password)
	if err != nil {
		log.Fatalf("failed to generate password hash: %v", err)
		return model.User{}, err
	}

	user := db.CreateUserParams{
		Email:        email,
		Username:     pgtype.Text{String: username, Valid: true},
		PasswordHash: passHash,
		Role:         db.RolesAdmin,
		Origin:       db.OriginsNative,
	}

	newAdmin, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Fatalf("failed to create admin: %v", err)
		return model.User{}, err
	}

	modelUser := model.DBUserToModelUser(newAdmin)[0]
	return modelUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		log.Fatalf("failed to delete user: %v", err)
		return err
	}

	return nil
}

func (s *Service) GetUser(ctx context.Context, id int64) (model.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found: %w", err)
		}
		log.Fatalf("failed to get user: %v", err)
		return model.User{}, err
	}

	modelUser := model.DBUserToModelUser(user)[0]
	return modelUser, nil
}

func (s *Service) GetUsersByRole(ctx context.Context, role model.Roles, limit, offset int32) ([]model.User, error) {
	var newLimit pgtype.Int4
	if limit <= 0 {
		newLimit.Valid = false
	} else {
		newLimit.Int32 = int32(limit)
		newLimit.Valid = true
	}

	users, err := s.repo.GetUsersByRole(ctx, db.GetUsersByRoleParams{Role: db.Roles(role), LimitParam: newLimit, Offset: offset})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		log.Fatalf("failed to get user role: %v", err)
		return nil, err
	}

	return model.DBUserToModelUser(users...), nil
}

func (s *Service) GetUsersLikeUsername(ctx context.Context, username string, limit, offset int32) ([]model.User, error) {
	wildcard := "%" + username + "%"

	var newLimit pgtype.Int4
	if limit <= 0 {
		newLimit.Valid = false
	} else {
		newLimit.Int32 = int32(limit)
		newLimit.Valid = true
	}

	users, err := s.repo.GetUsersLikeUsername(ctx, db.GetUsersLikeUsernameParams{Username: pgtype.Text{String: wildcard, Valid: true}, LimitParam: newLimit, Offset: offset})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		log.Fatalf("failed to get user like usersname: %v", err)
		return nil, err
	}

	return model.DBUserToModelUser(users...), nil
}

func (s *Service) GetuserByEmailOrUsername(ctx context.Context, email, username string) (model.User, error) {
	user, err := s.repo.GetuserByEmailOrUsername(ctx, db.GetuserByEmailOrUsernameParams{Email: email, Username: pgtype.Text{String: username, Valid: true}})
	if err != nil {
		return model.User{}, err
	}

	return model.DBUserToModelUser(user)[0], nil
}
