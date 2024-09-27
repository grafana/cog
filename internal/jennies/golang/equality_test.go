package golang_test

import (
	"testing"

	"github.com/grafana/cog/testdata/generated/equality"
	"github.com/stretchr/testify/require"
)

func TestEquality_Struct(t *testing.T) {
	req := require.New(t)

	req.True(equality.Container{StringField: "foo"}.Equals(equality.Container{StringField: "foo"}))
	req.True(equality.Container{IntField: 42}.Equals(equality.Container{IntField: 42}))
	req.True(equality.Container{EnumField: equality.DirectionBottom}.Equals(equality.Container{EnumField: equality.DirectionBottom}))
	req.True(equality.Container{RefField: equality.Variable{Name: "var"}}.Equals(equality.Container{RefField: equality.Variable{Name: "var"}}))

	req.False(equality.Container{StringField: "foo"}.Equals(equality.Container{StringField: "bar"}))
	req.False(equality.Container{IntField: 42}.Equals(equality.Container{IntField: 24}))
	req.False(equality.Container{EnumField: equality.DirectionBottom}.Equals(equality.Container{EnumField: equality.DirectionLeft}))
	req.False(equality.Container{RefField: equality.Variable{Name: "var"}}.Equals(equality.Container{RefField: equality.Variable{Name: "variable"}}))

	req.True(equality.Container{
		StringField: "foo",
		IntField:    42,
		EnumField:   equality.DirectionLeft,
		RefField:    equality.Variable{Name: "var"},
	}.Equals(equality.Container{
		StringField: "foo",
		IntField:    42,
		EnumField:   equality.DirectionLeft,
		RefField:    equality.Variable{Name: "var"},
	}))

	req.False(equality.Container{StringField: "foo"}.Equals(equality.Container{IntField: 42}))
	req.False(equality.Container{RefField: equality.Variable{Name: "var"}}.Equals(equality.Container{RefField: equality.Variable{Name: "other"}}))
}

func TestEquality_Struct_WithOptionalFields(t *testing.T) {
	req := require.New(t)

	// nil everywhere
	//nolint: gocritic
	req.True(equality.Optionals{}.Equals(equality.Optionals{}))

	req.True(equality.Optionals{
		StringField: toPtr("some string"),
	}.Equals(equality.Optionals{
		StringField: toPtr("some string"),
	}))
	req.True(equality.Optionals{
		EnumField: toPtr(equality.DirectionBottom),
	}.Equals(equality.Optionals{
		EnumField: toPtr[equality.Direction]("bottom"),
	}))
	req.True(equality.Optionals{
		RefField: &equality.Variable{Name: "foo"},
	}.Equals(equality.Optionals{
		RefField: &equality.Variable{Name: "foo"},
	}))

	req.False(equality.Optionals{
		StringField: toPtr("some string"),
	}.Equals(equality.Optionals{
		StringField: toPtr("other string"),
	}))
	req.False(equality.Optionals{
		EnumField: toPtr(equality.DirectionBottom),
	}.Equals(equality.Optionals{
		EnumField: toPtr[equality.Direction]("top"),
	}))
	req.False(equality.Optionals{
		RefField: &equality.Variable{Name: "foo"},
	}.Equals(equality.Optionals{
		RefField: &equality.Variable{Name: "bar"},
	}))

	req.False(equality.Optionals{
		StringField: toPtr("some string"),
	}.Equals(equality.Optionals{
		StringField: nil,
	}))
	req.False(equality.Optionals{
		StringField: nil,
	}.Equals(equality.Optionals{
		StringField: toPtr("some string"),
	}))
	req.False(equality.Optionals{
		EnumField: toPtr(equality.DirectionBottom),
	}.Equals(equality.Optionals{
		EnumField: nil,
	}))
	req.False(equality.Optionals{
		RefField: &equality.Variable{Name: "foo"},
	}.Equals(equality.Optionals{
		RefField: nil,
	}))
}

func TestEquality_Struct_WithArrays(t *testing.T) {
	req := require.New(t)

	// nil everywhere
	//nolint: gocritic
	req.True(equality.Arrays{}.Equals(equality.Arrays{}))

	req.True(equality.Arrays{Ints: []int64{}}.Equals(equality.Arrays{Ints: []int64{}}))
	req.True(equality.Arrays{Ints: []int64{3, 4, 5}}.Equals(equality.Arrays{Ints: []int64{3, 4, 5}}))

	req.True(equality.Arrays{Strings: []string{}}.Equals(equality.Arrays{Strings: []string{}}))
	req.True(equality.Arrays{Strings: []string{"foo", "bar", "baz"}}.Equals(equality.Arrays{Strings: []string{"foo", "bar", "baz"}}))

	req.True(equality.Arrays{ArrayOfArray: [][]string{}}.Equals(equality.Arrays{ArrayOfArray: [][]string{}}))
	req.True(equality.Arrays{ArrayOfArray: [][]string{{"foo", "bar"}}}.Equals(equality.Arrays{ArrayOfArray: [][]string{{"foo", "bar"}}}))
	req.True(equality.Arrays{ArrayOfArray: [][]string{{"foo", "bar"}, {"foo"}}}.Equals(equality.Arrays{ArrayOfArray: [][]string{{"foo", "bar"}, {"foo"}}}))

	req.True(equality.Arrays{Refs: []equality.Variable{}}.Equals(equality.Arrays{Refs: []equality.Variable{}}))
	req.True(equality.Arrays{Refs: []equality.Variable{{Name: "foo"}}}.Equals(equality.Arrays{Refs: []equality.Variable{{Name: "foo"}}}))

	req.True(equality.Arrays{ArrayOfAny: []any{}}.Equals(equality.Arrays{ArrayOfAny: []any{}}))
	req.True(equality.Arrays{ArrayOfAny: []any{"foo", 1}}.Equals(equality.Arrays{ArrayOfAny: []any{"foo", 1}}))

	req.False(equality.Arrays{Ints: []int64{1}}.Equals(equality.Arrays{Ints: []int64{}}))
	req.False(equality.Arrays{Ints: []int64{1, 2}}.Equals(equality.Arrays{Ints: []int64{3, 4}}))

	req.False(equality.Arrays{ArrayOfArray: [][]string{{"foo", "bar"}}}.Equals(equality.Arrays{ArrayOfArray: [][]string{{"foo"}}}))
	req.False(equality.Arrays{ArrayOfArray: [][]string{{"foo"}, {"bar"}}}.Equals(equality.Arrays{ArrayOfArray: [][]string{{"foo"}}}))
	req.False(equality.Arrays{ArrayOfArray: [][]string{{"foo"}, {"bar"}}}.Equals(equality.Arrays{ArrayOfArray: [][]string{{"foo"}, {"other"}}}))

	req.False(equality.Arrays{ArrayOfAny: []any{"foo", 1}}.Equals(equality.Arrays{ArrayOfAny: []any{"foo"}}))
	req.False(equality.Arrays{ArrayOfAny: []any{"bar"}}.Equals(equality.Arrays{ArrayOfAny: []any{"foo"}}))
}

func TestEquality_Struct_WithMaps(t *testing.T) {
	req := require.New(t)

	// nil everywhere
	//nolint: gocritic
	req.True(equality.Maps{}.Equals(equality.Maps{}))

	req.True(equality.Maps{
		Ints: map[string]int64{"foo": 42},
	}.Equals(equality.Maps{
		Ints: map[string]int64{"foo": 42},
	}))
	req.True(equality.Maps{
		Strings: map[string]string{"foo": "bar"},
	}.Equals(equality.Maps{
		Strings: map[string]string{"foo": "bar"},
	}))
	req.True(equality.Maps{
		Refs: map[string]equality.Variable{
			"foo": {Name: "foo"},
		},
	}.Equals(equality.Maps{
		Refs: map[string]equality.Variable{
			"foo": {Name: "foo"},
		},
	}))
	req.True(equality.Maps{
		StringToAny: map[string]any{
			"foo": 42,
			"bar": true,
		},
	}.Equals(equality.Maps{
		StringToAny: map[string]any{
			"foo": 42,
			"bar": true,
		},
	}))

	req.False(equality.Maps{
		Ints: map[string]int64{"foo": 42},
	}.Equals(equality.Maps{
		Ints: map[string]int64{"foo": 1},
	}))
	req.False(equality.Maps{
		Ints: map[string]int64{"foo": 42, "bar": 24},
	}.Equals(equality.Maps{
		Ints: map[string]int64{"foo": 42},
	}))

	req.False(equality.Maps{
		Strings: map[string]string{"foo": "foo"},
	}.Equals(equality.Maps{
		Strings: map[string]string{"foo": "not foo"},
	}))
	req.False(equality.Maps{
		Strings: map[string]string{"foo": "foo", "bar": "bar"},
	}.Equals(equality.Maps{
		Strings: map[string]string{"foo": "foo"},
	}))

	req.False(equality.Maps{
		Refs: map[string]equality.Variable{
			"foo": {Name: "foo"},
		},
	}.Equals(equality.Maps{
		Refs: map[string]equality.Variable{
			"foo": {Name: "bar"},
		},
	}))
	req.False(equality.Maps{
		Refs: map[string]equality.Variable{
			"foo": {Name: "foo"},
		},
	}.Equals(equality.Maps{
		Refs: map[string]equality.Variable{
			"foo": {Name: "foo"},
			"bar": {Name: "bar"},
		},
	}))
}

func toPtr[V any](input V) *V {
	return &input
}
