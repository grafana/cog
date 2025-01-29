package languages

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestContext_ResolveAsBuilder(t *testing.T) {
	fooObj := ast.NewObject("foo", "Foo", ast.NewStruct(ast.NewStructField("bar", ast.String())))
	bizObj := ast.NewObject("foo", "Biz", ast.NewStruct(ast.NewStructField("bar", ast.String())))
	fooOrBiz := ast.NewObject("foo", "FooOrBiz", ast.NewDisjunction([]ast.Type{
		ast.NewRef("foo", "Foo"),
		ast.NewRef("foo", "Biz"),
	}))

	context := Context{
		Schemas: []*ast.Schema{
			{
				Package: "foo",
				Objects: testutils.ObjectsMap(
					fooObj,
					bizObj,
					fooOrBiz,
					ast.NewObject("foo", "Bar", ast.NewStruct(ast.NewStructField("bar", ast.String()))),
				),
			},
		},
		Builders: []ast.Builder{
			{
				Name: "Foo",
				For:  fooObj,
			},
			{
				Name: "Biz",
				For:  bizObj,
			},
		},
	}

	testCases := []struct {
		desc           string
		input          ast.Type
		expectedResult bool
	}{
		{
			desc:           "ref to buildable",
			input:          ast.NewRef("foo", "Foo"),
			expectedResult: true,
		},
		{
			desc:           "ref to NOT buildable",
			input:          ast.NewRef("foo", "Bar"),
			expectedResult: false,
		},

		{
			desc:           "array of ref to buildable",
			input:          ast.NewArray(ast.NewRef("foo", "Foo")),
			expectedResult: true,
		},
		{
			desc:           "array of to ref to NOT buildable",
			input:          ast.NewArray(ast.NewRef("foo", "Bar")),
			expectedResult: false,
		},

		{
			desc:           "map of string to ref to buildable",
			input:          ast.NewMap(ast.String(), ast.NewRef("foo", "Foo")),
			expectedResult: true,
		},
		{
			desc:           "map of string to ref to NOT buildable",
			input:          ast.NewMap(ast.String(), ast.NewRef("foo", "Bar")),
			expectedResult: false,
		},

		{
			desc:           "disjunction including ref to buildable",
			input:          ast.NewDisjunction([]ast.Type{ast.String(), ast.NewRef("foo", "Foo")}),
			expectedResult: true,
		},
		{
			desc:           "disjunction of NOT buildable types",
			input:          ast.NewDisjunction(ast.Types{ast.String(), ast.NewRef("foo", "Bar")}),
			expectedResult: false,
		},

		{
			desc:           "ref to disjunction of buildable",
			input:          ast.NewRef("foo", "FooOrBiz"),
			expectedResult: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			req := require.New(t)
			req.Equal(testCase.expectedResult, context.ResolveToBuilder(testCase.input))
		})
	}
}
