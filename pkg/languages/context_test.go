package languages

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
	"github.com/stretchr/testify/require"
)

func TestContext_ResolveAsBuilder(t *testing.T) {
	fooObj := ir.NewObject("foo", "Foo", ir.NewStruct(ir.NewStructField("bar", ir.String())))
	bizObj := ir.NewObject("foo", "Biz", ir.NewStruct(ir.NewStructField("bar", ir.String())))
	fooOrBiz := ir.NewObject("foo", "FooOrBiz", ir.NewDisjunction([]ir.Type{
		ir.NewRef("foo", "Foo"),
		ir.NewRef("foo", "Biz"),
	}))

	context := Context{
		Schemas: []*ir.Schema{
			{
				Package: "foo",
				Objects: testutils.ObjectsMap(
					fooObj,
					bizObj,
					fooOrBiz,
					ir.NewObject("foo", "Bar", ir.NewStruct(ir.NewStructField("bar", ir.String()))),
				),
			},
		},
		Builders: []ir.Builder{
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
		input          ir.Type
		expectedResult bool
	}{
		{
			desc:           "ref to buildable",
			input:          ir.NewRef("foo", "Foo"),
			expectedResult: true,
		},
		{
			desc:           "ref to NOT buildable",
			input:          ir.NewRef("foo", "Bar"),
			expectedResult: false,
		},

		{
			desc:           "array of ref to buildable",
			input:          ir.NewArray(ir.NewRef("foo", "Foo")),
			expectedResult: true,
		},
		{
			desc:           "array of to ref to NOT buildable",
			input:          ir.NewArray(ir.NewRef("foo", "Bar")),
			expectedResult: false,
		},

		{
			desc:           "map of string to ref to buildable",
			input:          ir.NewMap(ir.String(), ir.NewRef("foo", "Foo")),
			expectedResult: true,
		},
		{
			desc:           "map of string to ref to NOT buildable",
			input:          ir.NewMap(ir.String(), ir.NewRef("foo", "Bar")),
			expectedResult: false,
		},

		{
			desc:           "disjunction including ref to buildable",
			input:          ir.NewDisjunction([]ir.Type{ir.String(), ir.NewRef("foo", "Foo")}),
			expectedResult: true,
		},
		{
			desc:           "disjunction of NOT buildable types",
			input:          ir.NewDisjunction(ir.Types{ir.String(), ir.NewRef("foo", "Bar")}),
			expectedResult: false,
		},

		{
			desc:           "ref to disjunction of buildable",
			input:          ir.NewRef("foo", "FooOrBiz"),
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
