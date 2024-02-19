package model

import (
	"time"

	db "github.com/izzanzahrial/skeleton/db/sqlc"
)

type Roles string

const (
	RolesAdmin Roles = "admin"
	RolesUser  Roles = "user"
)

type User struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash []byte    `json:"password_hash"`
	Role         Roles     `json:"role"`
}

// DBUserToModelUser converts a DB user to a model user
func DBUserToModelUser(users ...db.User) []User {
	var modelUsers []User

	for _, u := range users {
		modelUsers = append(modelUsers, User{
			ID:           u.ID,
			CreatedAt:    u.CreatedAt.Time,
			UpdatedAt:    u.UpdatedAt.Time,
			DeletedAt:    u.DeletedAt.Time,
			Email:        u.Email,
			Username:     u.Username,
			PasswordHash: u.PasswordHash,
			Role:         Roles(u.Role),
		})
	}

	return modelUsers
}
