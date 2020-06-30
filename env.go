package appconfig

import (
	"context"
	"os"
	"strings"

	"github.com/heetch/confita/backend"
)

func NewEnvBackend() backend.Backend {
	return backend.Func("env", func(ctx context.Context, key string) ([]byte, error) {
		if val := os.Getenv(key); val != "" {
			return []byte(val), nil
		}
		key = strings.Replace(strings.ToUpper(key), ".", "_", -1)
		key = strings.Replace(key, "-", "_", -1)

		if val := os.Getenv(key); val != "" {
			return []byte(val), nil
		}
		return nil, backend.ErrNotFound
	})
}
