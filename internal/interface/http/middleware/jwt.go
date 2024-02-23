package middleware

import (
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/izzanzahrial/skeleton/internal/model"
	"github.com/izzanzahrial/skeleton/pkg/token"
)

// using authorization headers but from echo
func IsAuthorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tkn := c.Get("user").(*jwt.Token)
		claims, ok := tkn.Claims.(*token.JwtCustomClaims)
		if !ok && !tkn.Valid {
			return echo.ErrUnauthorized
		}

		if claims.Role != model.RolesAdmin {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

// using authorization header
func IsAuthorizeHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authtoken := c.Request().Header.Get("Authorization")
		if authtoken == "" {
			return echo.ErrUnauthorized
		}

		tokenString := strings.Split(authtoken, " ")[1]
		tkn, err := jwt.ParseWithClaims(tokenString, &token.JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			// TODO: move secret key somewhere else
			return []byte("secret"), nil
		})
		if err != nil {
			log.Fatalf("failed to parse token: %v", err)
			return echo.ErrUnauthorized
		}

		claims, ok := tkn.Claims.(*token.JwtCustomClaims)
		if !ok && !tkn.Valid {
			return echo.ErrUnauthorized
		}

		if claims.Role != model.RolesAdmin {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

func IsAuthenticated() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		// TODO: move secret key somewhere else
		SigningKey: []byte("secret"),
	}

	return echojwt.WithConfig(config)
}
