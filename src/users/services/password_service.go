package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/terdia/greenlight/src/users/entities"
)

const (
	cost = 12
)

type PasswordHashService interface {
	Hash(password *entities.Password) error
	Verify(hash []byte, plainText string) (bool, error)
}

type bcryptPasswordService struct {
}

func NewPasswordService() PasswordHashService {
	return &bcryptPasswordService{}
}

func (service *bcryptPasswordService) Hash(password *entities.Password) error {
	if password.PlainText == nil {
		return errors.New("plaintext password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password.PlainText), cost)
	if err != nil {
		return err
	}

	password.Hash = hashedPassword

	return nil
}

func (service *bcryptPasswordService) Verify(hash []byte, plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
