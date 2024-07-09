package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func Parse[T any](token string) (T, error) {
	var result T

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return result, fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return result, fmt.Errorf("failed to decode token payload: %v", err)
	}

	err = json.Unmarshal(payload, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal token payload: %v", err)
	}

	var payloadMap map[string]interface{}
	err = json.Unmarshal(payload, &payloadMap)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal token payload for expiration check: %v", err)
	}

	if exp, ok := payloadMap["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return result, fmt.Errorf("token has expired")
		}
	}

	return result, nil
}
