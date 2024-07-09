package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestPayload struct {
	Exp int64  `json:"exp"`
	Foo string `json:"foo"`
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    TestPayload
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   createTestToken(TestPayload{Exp: time.Now().Add(time.Hour).Unix(), Foo: "bar"}),
			want:    TestPayload{Exp: time.Now().Add(time.Hour).Unix(), Foo: "bar"},
			wantErr: false,
		},
		{
			name:    "invalid token format parts > 3",
			token:   "invalid.token.format.signature",
			want:    TestPayload{},
			wantErr: true,
		},
		{
			name:    "invalid token format",
			token:   "invalid.token.format",
			want:    TestPayload{},
			wantErr: true,
		},
		{
			name:    "expired token",
			token:   createTestToken(TestPayload{Exp: time.Now().Add(-time.Hour).Unix(), Foo: "bar"}),
			want:    TestPayload{},
			wantErr: true,
		},
		{
			name:    "invalid payload",
			token:   createInvalidPayloadToken(),
			want:    TestPayload{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestPayload
			got, err := Parse[TestPayload](tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func createTestToken(payload TestPayload) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, _ := json.Marshal(payload)
	payloadStr := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature := "signature"

	return fmt.Sprintf("%s.%s.%s", header, payloadStr, signature)
}

func createInvalidPayloadToken() string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadStr := "invalidPayloadString"
	signature := "signature"

	return fmt.Sprintf("%s.%s.%s", header, payloadStr, signature)
}
