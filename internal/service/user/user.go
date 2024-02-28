package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type Service struct {
	repo userRepo
	slog *slog.Logger
}

func NewService(repo userRepo, slog *slog.Logger) *Service {
	return &Service{
		repo: repo,
		slog: slog,
	}
}

func (s *Service) CreateUser(ctx context.Context, email, username, password string) (model.User, error) {
	passHash, err := pass.Generate(password)
	if err != nil {
		s.slog.Error("failed to generate password hash", slog.String("error", err.Error()))
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
		s.slog.Error("failed to create user", slog.String("error", err.Error()))
		return model.User{}, err
	}

	modelUser := model.DBUserToModelUser(newUser)[0]
	return modelUser, nil
}

func (s *Service) CreateAdmin(ctx context.Context, email, username, password string) (model.User, error) {
	passHash, err := pass.Generate(password)
	if err != nil {
		s.slog.Error("failed to generate password hash", slog.String("error", err.Error()))
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
		s.slog.Error("failed to create user", slog.String("error", err.Error()))
		return model.User{}, err
	}

	modelUser := model.DBUserToModelUser(newAdmin)[0]
	return modelUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		s.slog.Error("failed to delete user", slog.String("error", err.Error()))
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
		s.slog.Error("failed to get user", slog.String("error", err.Error()))
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
		s.slog.Error("failed to get users by role", slog.String("error", err.Error()))
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
		s.slog.Error("failed to get users using like username", slog.String("error", err.Error()))
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

func (s *Service) UpdateUser(ctx context.Context, id int64, email, username, password *string) (model.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.slog.Error("failed to get user", slog.String("error", err.Error()))
			return model.User{}, fmt.Errorf("user not found: %w", err)
		}
		return model.User{}, err
	}

	if email != nil {
		user.Email = *email
	}
	if username != nil {
		user.Username = pgtype.Text{String: *username, Valid: true}
	}
	if password != nil {
		passHash, err := pass.Generate(*password)
		if err != nil {
			s.slog.Error("failed to generate password hash", slog.String("error", err.Error()))
			return model.User{}, err
		}
		user.PasswordHash = passHash
	}

	updatedUser, err := s.repo.UpdateUser(ctx, db.UpdateUserParams{Email: user.Email, Username: user.Username, PasswordHash: user.PasswordHash, ID: id})
	if err != nil {
		s.slog.Error("failed to update user", slog.String("error", err.Error()))
		return model.User{}, err
	}

	return model.DBUserToModelUser(updatedUser)[0], nil
}
