package paseto

import (
	"fmt"
	"github.com/go-lite/lite/tokens"
	"github.com/o1egl/paseto"
)

func Parse[T tokens.Claims](token string, key []byte) (T, error) {
	var val T
	var footer string

	v2 := paseto.NewV2()

	err := v2.Decrypt(token, key, &val, &footer)
	if err != nil {
		return val, fmt.Errorf("failed to decrypt token: %w", err)
	}

	if !val.Valid() {
		return val, fmt.Errorf("token is invalid")
	}

	return val, nil
}
