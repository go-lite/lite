package paseto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

var (
	jsonMarshal   = json.Marshal
	jsonUnmarshal = json.Unmarshal
)

func Parse[T any](token string, key []byte) (T, error) {
	var result T
	var jsonToken paseto.JSONToken
	var footer string

	v2 := paseto.NewV2()

	err := v2.Decrypt(token, key, &jsonToken, &footer)
	if err != nil {
		return result, fmt.Errorf("failed to decrypt token: %w", err)
	}

	if jsonToken.Expiration.Before(time.Now()) {
		return result, fmt.Errorf("token has expired")
	}

	payload, err := jsonMarshal(jsonToken)
	if err != nil {
		return result, fmt.Errorf("failed to marshal token payload: %w", err)
	}

	err = jsonUnmarshal(payload, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal token payload: %w", err)
	}

	return result, nil
}
