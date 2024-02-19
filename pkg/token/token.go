package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/izzanzahrial/skeleton/internal/model"
)

type Token struct {
	// user model is the main identifier
	User      *model.User
	PlainText string
	Hash      []byte
	Expiry    time.Duration
}

func New(userID int64, ttl time.Duration) (*Token, error) {
	token := &Token{
		User:   &model.User{ID: userID},
		Expiry: ttl,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %v", err)
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}
