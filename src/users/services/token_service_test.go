package services

import (
	"testing"
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
)

func TestCreateToken(t *testing.T) {

	tests := []struct {
		name       string
		userId     custom_type.ID
		ttl        time.Duration
		scope      string
		wantScope  string
		wantLength int
	}{
		{"Generate Activation token", custom_type.ID(4), 5 * time.Minute, data.TokenScopeActivation, "activation", 32},
		{"Generate Authentication token", custom_type.ID(5), 5 * time.Minute, data.TokenScopeAuthentication, "authentication", 32},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := generateToken(test.userId, test.ttl, test.scope)

			if err != nil {
				t.Errorf("want error to be %v; got %s", nil, err.Error())
			}

			if token.Scope != test.scope {
				t.Errorf("want %v; got %s", test.wantScope, token.Scope)
			}

			if len(token.Hash) != test.wantLength {
				t.Errorf("want %v; got %d", test.wantLength, len(token.Hash))
			}

			if !token.Expiry.After(time.Now()) {
				t.Errorf("want %v in the future; got %s", test.ttl, token.Expiry)
			}
		})
	}
}
