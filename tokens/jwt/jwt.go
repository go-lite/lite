package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var jsonUnmarshal = json.Unmarshal

func Parse[T any](token string) (T, error) {
	var result T

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return result, errors.New("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return result, fmt.Errorf("failed to decode token payload: %w", err)
	}

	err = json.Unmarshal(payload, &result)
	if err != nil {
		return result, fmt.Errorf(
			"failed to unmarshal token payload: %w",
			err,
		)
	}

	var payloadMap map[string]interface{}

	err = jsonUnmarshal(payload, &payloadMap)
	if err != nil {
		return result, fmt.Errorf(
			"failed to unmarshal token payload for expiration check: %w",
			err,
		)
	}

	if exp, ok := payloadMap["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return result, errors.New("token has expired")
		}
	}

	return result, nil
}
