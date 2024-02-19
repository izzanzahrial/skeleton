package user

import (
	"context"

	db "github.com/izzanzahrial/skeleton/db/sqlc"
	"github.com/izzanzahrial/skeleton/internal/model"
	pass "github.com/izzanzahrial/skeleton/pkg/password"
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
		return model.User{}, err
	}

	user := db.CreateUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: passHash,
		Role:         db.RolesUser,
	}

	newUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	modelUser := model.DBUserToModelUser(newUser)[0]
	return modelUser, nil
}

func (s *Service) CreateAdmin(ctx context.Context, email, username, password string) (model.User, error) {
	passHash, err := pass.Generate(password)
	if err != nil {
		return model.User{}, err
	}

	user := db.CreateUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: passHash,
		Role:         db.RolesAdmin,
	}

	newUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	modelUser := model.DBUserToModelUser(newUser)[0]
	return modelUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *Service) GetUser(ctx context.Context, id int64) (model.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
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

	users, err := s.repo.GetUsersByRole(ctx, db.GetUsersByRoleParams{Role: db.Roles(role), LimitArg: newLimit, Offset: offset})
	if err != nil {
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

	users, err := s.repo.GetUsersLikeUsername(ctx, db.GetUsersLikeUsernameParams{Username: wildcard, LimitArg: newLimit, Offset: offset})
	if err != nil {
		return nil, err
	}

	return model.DBUserToModelUser(users...), nil
}

func (s *Service) GetuserByEmailOrUsername(ctx context.Context, email, username string) (model.User, error) {
	user, err := s.repo.GetuserByEmailOrUsername(ctx, db.GetuserByEmailOrUsernameParams{Email: email, Username: username})
	if err != nil {
		return model.User{}, err
	}

	return model.DBUserToModelUser(user)[0], nil
}
