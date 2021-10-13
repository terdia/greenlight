package repositories

import (
	"github.com/terdia/greenlight/src/users/entities"
)

type TokenRepository interface {
	Create(token *entities.Token) error
	DeleteAllForUserByScope(scope string, userID int64) error
}
