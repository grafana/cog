package httputil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidURL(t *testing.T) {
	t.Run("invalid schema", func(t *testing.T) {
		url := "http://localhost:8080"
		_, err := LoadURL(context.Background(), url)
		require.ErrorContains(t, err, "unsupported scheme 'http'")
	})

	t.Run("invalid host", func(t *testing.T) {
		url := "https://unknown.host:8080"
		_, err := LoadURL(context.Background(), url)
		require.ErrorContains(t, err, "unsupported host 'unknown.host'")
	})
}
