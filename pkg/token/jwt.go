package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/izzanzahrial/skeleton/internal/model"
)

type JwtCustomClaims struct {
	UserID int64       `json:"user_id"`
	Role   model.Roles `json:"role"`
	jwt.RegisteredClaims
}

func NewJWT(userID int64, role model.Roles) (string, error) {
	expiry := time.Now().Add(time.Hour * 24)

	claims := &JwtCustomClaims{
		userID,
		role,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expiry)},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// TODO: create secret keys for jwt
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt token: %v", err)
	}

	return t, nil
}
