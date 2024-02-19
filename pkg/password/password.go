package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Check(plainTextPass string, passHash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(passHash, []byte(plainTextPass))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, errors.New("failed password does not match")
		default:
			return false, err
		}
	}

	return true, nil
}

func Generate(plaintTextPass string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(plaintTextPass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return passwordHash, nil
}
