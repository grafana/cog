package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypes_HasOnlyScalarOrArray(t *testing.T) {
	testCases := []struct {
		description string
		types       Types
		expected    bool
	}{
		{
			description: "only scalars",
			types: []Type{
				NewScalar(KindString),
				NewScalar(KindBool),
			},
			expected: true,
		},
		{
			description: "scalars and an array of scalars",
			types: []Type{
				NewScalar(KindString),
				NewArray(NewScalar(KindInt8)),
			},
			expected: true,
		},
		{
			description: "scalars and an array of refs",
			types: []Type{
				NewScalar(KindString),
				NewArray(NewRef("pkg", "SomeType")),
			},
			expected: true,
		},
		{
			description: "ref",
			types: []Type{
				NewRef("pkg", "SomeType"),
			},
			expected: false,
		},
		{
			description: "scalars and ref",
			types: []Type{
				NewScalar(KindString),
				NewRef("pkg", "SomeType"),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.description, func(t *testing.T) {
			req := require.New(t)

			req.Equal(tc.expected, tc.types.HasOnlyScalarOrArrayOrMap())
		})
	}
}

func TestTypes_HasOnlyRefs(t *testing.T) {
	testCases := []struct {
		description string
		types       Types
		expected    bool
	}{
		{
			description: "only scalars",
			types: []Type{
				NewScalar(KindString),
				NewScalar(KindBool),
			},
			expected: false,
		},
		{
			description: "scalars and an array of scalars",
			types: []Type{
				NewScalar(KindString),
				NewArray(NewScalar(KindInt8)),
			},
			expected: false,
		},
		{
			description: "refs",
			types: []Type{
				NewRef("pkg", "SomeType"),
				NewRef("pkg", "SomeOtherType"),
			},
			expected: true,
		},
		{
			description: "ref",
			types: []Type{
				NewRef("pkg", "SomeType"),
			},
			expected: true,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.description, func(t *testing.T) {
			req := require.New(t)

			req.Equal(tc.expected, tc.types.HasOnlyRefs())
		})
	}
}

func TestArrayType_IsArrayOfScalars(t *testing.T) {
	testCases := []struct {
		description string
		array       ArrayType
		Expected    bool
	}{
		{
			description: "array of scalars",
			array: ArrayType{
				ValueType: NewScalar(KindString),
			},
			Expected: true,
		},
		{
			description: "array of array of scalars",
			array: ArrayType{
				ValueType: NewArray(NewScalar(KindString)),
			},
			Expected: true,
		},
		{
			description: "array of refs",
			array: ArrayType{
				ValueType: NewRef("pkg", "SomeType"),
			},
			Expected: false,
		},
		{
			description: "array of structs",
			array: ArrayType{
				ValueType: NewStruct(NewStructField("Foo", NewScalar(KindString))),
			},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.description, func(t *testing.T) {
			req := require.New(t)

			req.Equal(tc.Expected, tc.array.IsArrayOfScalars())
		})
	}
}
