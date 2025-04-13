package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hashed(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash: %w", err)
	}
	return string(hash), nil
}

func Compare(hashed, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashed),
		[]byte(password),
	); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			return false, fmt.Errorf("mismatch password: %w", err)
		}
		return false, err
	}
	return true, nil
}
