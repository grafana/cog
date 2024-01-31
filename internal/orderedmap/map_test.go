package orderedmap

import (
	"encoding/json"
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

func TestMap_MarshalJSON(t *testing.T) {
	req := require.New(t)
	orderedMap := New[string, string]()

	marshaled, err := json.Marshal(orderedMap)
	req.NoError(err)
	req.Equal(`{}`, string(marshaled))

	orderedMap.Set("foo", "first")
	orderedMap.Set("bar", "second")

	marshaled, err = json.Marshal(orderedMap)
	req.NoError(err)
	req.Equal(`{"foo":"first","bar":"second"}`, string(marshaled))
}

func TestMap_UnmarshalJSON_empty(t *testing.T) {
	req := require.New(t)
	orderedMap := New[string, string]()

	req.NoError(json.Unmarshal([]byte(`{}`), orderedMap))
	req.Equal(0, orderedMap.Len())
}

func TestMap_UnmarshalJSON_string_to_string(t *testing.T) {
	req := require.New(t)
	orderedMap := New[string, string]()

	req.NoError(json.Unmarshal([]byte(`{"foo":"first","bar":"second","aaa":"third"}`), orderedMap))
	req.Equal(3, orderedMap.Len())
	req.Equal("first", orderedMap.At(0))
	req.Equal("second", orderedMap.At(1))
	req.Equal("third", orderedMap.At(2))
}

func TestMap_UnmarshalJSON_string_to_struct(t *testing.T) {
	req := require.New(t)
	orderedMap := New[string, struct {
		Title string `json:"title"`
	}]()

	req.NoError(json.Unmarshal([]byte(`{"foo":{"title":"first"},"bar":{"title":"second"},"aaa":{"title":"third"}}`), orderedMap))
	req.Equal(3, orderedMap.Len())
	req.Equal("first", orderedMap.At(0).Title)
	req.Equal("second", orderedMap.At(1).Title)
	req.Equal("third", orderedMap.At(2).Title)
}
