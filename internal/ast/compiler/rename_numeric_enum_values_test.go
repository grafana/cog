package compiler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func TestRenameNumericEnumValues(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("pkg", "NotAnEnumStruct", ast.String()),

		ast.NewObject("pkg", "AnEnumWithNumericValues", ast.NewEnum([]ast.EnumValue{
			{
				Type:  ast.NewScalar(ast.KindInt64),
				Name:  "1",
				Value: 1,
			},
			{
				Type:  ast.NewScalar(ast.KindInt64),
				Name:  "2",
				Value: 2,
			},
		})),

		ast.NewObject("pkg", "AnEnumWithNoNumericValues", ast.NewEnum([]ast.EnumValue{
			{
				Type:  ast.String(),
				Name:  "Hide",
				Value: "hide",
			},
			{
				Type:  ast.String(),
				Name:  "DontHide",
				Value: "dont_hide",
			},
		})),
	}

	// Prepare expected output
	expected := []ast.Object{
		ast.NewObject("pkg", "NotAnEnumStruct", ast.String()),

		ast.NewObject("pkg", "AnEnumWithNumericValues", ast.NewEnum([]ast.EnumValue{
			{
				Type:  ast.NewScalar(ast.KindInt64),
				Name:  "N1",
				Value: 1,
			},
			{
				Type:  ast.NewScalar(ast.KindInt64),
				Name:  "N2",
				Value: 2,
			},
		})),

		ast.NewObject("pkg", "AnEnumWithNoNumericValues", ast.NewEnum([]ast.EnumValue{
			{
				Type:  ast.String(),
				Name:  "Hide",
				Value: "hide",
			},
			{
				Type:  ast.String(),
				Name:  "DontHide",
				Value: "dont_hide",
			},
		})),
	}

	// Run the compiler pass
	runRenameNumericEnumValuesPass(t, objects, expected)
}

func runRenameNumericEnumValuesPass(t *testing.T, input []ast.Object, expectedOutput []ast.Object) {
	t.Helper()

	req := require.New(t)

	compilerPass := &RenameNumericEnumValues{}
	processedFiles, err := compilerPass.Process([]*ast.Schema{
		{
			Package: "test",
			Objects: input,
		},
	})
	req.NoError(err)
	req.Len(processedFiles, 1)
	req.Empty(cmp.Diff(expectedOutput, processedFiles[0].Objects))
}
