package services

import (
	"testing"

	"github.com/terdia/greenlight/src/users/entities"
)

func TestPasswordService(t *testing.T) {

	plainText := "testing"

	password := entities.Password{
		PlainText: nil,
	}

	passwordService := NewPasswordService()

	t.Run("nil password", func(t *testing.T) {
		err := passwordService.Hash(&password)

		if err.Error() != "plaintext password is required" {
			t.Errorf("want %s; got %v", "error: plaintext password is required", nil)
		}
	})

	t.Run("hash a plain text password", func(t *testing.T) {
		password.PlainText = &plainText
		err := passwordService.Hash(&password)

		if err != nil {
			t.Error(err.Error())
		}
	})

	t.Run("verify password hash ok", func(t *testing.T) {
		password.PlainText = &plainText
		ok, err := passwordService.Verify(password.Hash, plainText)

		if err != nil {
			t.Error(err.Error())
		}

		if ok != true {
			t.Errorf("want %v; got %v", true, false)
		}
	})

	t.Run("test password mismatch error", func(t *testing.T) {
		password.PlainText = &plainText
		ok, err := passwordService.Verify(password.Hash, "somezehhe")

		if err != nil {
			t.Error(err.Error())
		}

		if ok != false {
			t.Errorf("want %v; got %v", false, true)
		}
	})

}
