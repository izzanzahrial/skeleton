// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Origins string

const (
	OriginsNative Origins = "native"
	OriginsGoogle Origins = "google"
)

func (e *Origins) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Origins(s)
	case string:
		*e = Origins(s)
	default:
		return fmt.Errorf("unsupported scan type for Origins: %T", src)
	}
	return nil
}

type NullOrigins struct {
	Origins Origins `json:"origins"`
	Valid   bool    `json:"valid"` // Valid is true if Origins is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrigins) Scan(value interface{}) error {
	if value == nil {
		ns.Origins, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Origins.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrigins) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Origins), nil
}

type Roles string

const (
	RolesAdmin Roles = "admin"
	RolesUser  Roles = "user"
)

func (e *Roles) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Roles(s)
	case string:
		*e = Roles(s)
	default:
		return fmt.Errorf("unsupported scan type for Roles: %T", src)
	}
	return nil
}

type NullRoles struct {
	Roles Roles `json:"roles"`
	Valid bool  `json:"valid"` // Valid is true if Roles is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRoles) Scan(value interface{}) error {
	if value == nil {
		ns.Roles, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Roles.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRoles) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Roles), nil
}

type Post struct {
	ID        int64              `json:"id"`
	UserID    int64              `json:"user_id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
}

type User struct {
	ID           int64              `json:"id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	DeletedAt    pgtype.Timestamptz `json:"deleted_at"`
	Email        string             `json:"email"`
	Username     pgtype.Text        `json:"username"`
	PasswordHash []byte             `json:"password_hash"`
	Role         Roles              `json:"role"`
	FirstName    pgtype.Text        `json:"first_name"`
	LastName     pgtype.Text        `json:"last_name"`
	PictureUrl   pgtype.Text        `json:"picture_url"`
	RefreshToken pgtype.Text        `json:"refresh_token"`
	Origin       Origins            `json:"origin"`
}
