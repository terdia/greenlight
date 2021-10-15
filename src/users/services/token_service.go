package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/src/users/entities"
	"github.com/terdia/greenlight/src/users/repositories"
)

type TokenService interface {
	CreateNew(userId custom_type.ID, ttl time.Duration, scope string) (*entities.Token, error)
	DeleteByUserIdAndScope(userId custom_type.ID, scope string) error
}

type tokenService struct {
	repo repositories.TokenRepository
}

func NewTokenService(tokenRepository repositories.TokenRepository) TokenService {
	return &tokenService{repo: tokenRepository}
}

func (tsrv tokenService) CreateNew(userId custom_type.ID, ttl time.Duration, scope string) (*entities.Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = tsrv.repo.Create(token)

	return token, err
}

func (tsrv tokenService) DeleteByUserIdAndScope(userId custom_type.ID, scope string) error {
	return tsrv.repo.DeleteAllForUserByScope(scope, userId)
}

func generateToken(userId custom_type.ID, ttl time.Duration, scope string) (*entities.Token, error) {

	token := &entities.Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	// fill the byte slice with random bytes from your os CSPRNG.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))

	//convert it to a slice using the [:] operator
	token.Hash = hash[:]

	return token, nil
}
