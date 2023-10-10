package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap_Basic(t *testing.T) {
	req := require.New(t)
	orderedMap := New[string, string]()

	req.Equal(0, orderedMap.Len())
	req.False(orderedMap.Has("some-key"))

	orderedMap.Set("first-key", "first-value")
	orderedMap.Set("second-key", "second-value")
	orderedMap.Set("third-key", "third-value")

	req.Equal(3, orderedMap.Len())
	req.Equal("second-value", orderedMap.Get("second-key"))
	req.False(orderedMap.Has("unknown-key"))
	req.Equal("", orderedMap.Get("unknown-key"))

	iteratedKeyOrder := make([]string, 0, orderedMap.Len())
	orderedMap.Iterate(func(key string, _ string) {
		iteratedKeyOrder = append(iteratedKeyOrder, key)
	})

	req.Equal([]string{"first-key", "second-key", "third-key"}, iteratedKeyOrder)
}
