package paseto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/o1egl/paseto"
)

type TestPayload struct {
	Exp time.Time `json:"exp"`
	Foo string    `json:"foo"`
}

func (p TestPayload) Valid() bool {
	return p.Exp.After(time.Now())
}

func TestParse(t *testing.T) {
	secretKey := []byte("YELLOW SUBMARINE, BLACK WIZARDRY") // Utilisez une clé appropriée pour vos tests

	tests := []struct {
		name    string
		token   string
		key     []byte
		want    TestPayload
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   createTestToken(TestPayload{Exp: time.Now().Add(time.Hour), Foo: "bar"}, secretKey),
			key:     secretKey,
			want:    TestPayload{Exp: time.Now().Add(time.Hour), Foo: "bar"},
			wantErr: false,
		},
		{
			name:    "invalid token format",
			token:   "invalid.token.format",
			key:     secretKey,
			want:    TestPayload{},
			wantErr: true,
		},
		{
			name:    "expired token",
			token:   createTestToken(TestPayload{Exp: time.Now().Add(-time.Hour), Foo: "bar"}, secretKey),
			key:     secretKey,
			want:    TestPayload{},
			wantErr: true,
		},
		{
			name:    "invalid key",
			token:   createTestToken(TestPayload{Exp: time.Now().Add(time.Hour), Foo: "bar"}, secretKey),
			key:     []byte("INVALID_KEY"),
			want:    TestPayload{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestPayload
			got, err := Parse[TestPayload](tt.token, tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Foo, got.Foo)
				assert.WithinDuration(t, tt.want.Exp, got.Exp, time.Second)
			}
		})
	}
}

func createTestToken(payload TestPayload, key []byte) string {
	jsonToken := paseto.JSONToken{
		Expiration: payload.Exp,
	}
	jsonToken.Set("foo", payload.Foo)

	v2 := paseto.NewV2()
	token, err := v2.Encrypt(key, jsonToken, "")
	if err != nil {
		panic(fmt.Sprintf("failed to create test token: %v", err))
	}

	return token
}
