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

type Origins string

const (
	NativeOrigin Origins = "native"
	GoogleOrigin Origins = "google"
)

type User struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash []byte    `json:"password_hash"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PictureUrl   string    `json:"picture_url"`
	RefreshToken string    `json:"refresh_token"`
	Role         Roles     `json:"role"`
	Origin       Origins   `json:"origin"`
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
			Username:     u.Username.String,
			PasswordHash: u.PasswordHash,
			FirstName:    u.FirstName.String,
			LastName:     u.LastName.String,
			PictureUrl:   u.PictureUrl.String,
			RefreshToken: u.RefreshToken.String,
			Role:         Roles(u.Role),
			Origin:       Origins(u.Origin),
		})
	}

	return modelUsers
}
