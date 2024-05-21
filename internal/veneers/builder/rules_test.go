package builder

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestDuplicate(t *testing.T) {
	req := require.New(t)

	originalObject := ast.NewObject("pkg", "Dashboard", ast.NewStruct(
		ast.NewStructField("name", ast.String()),
	))
	argument := ast.Argument{Name: "title", Type: ast.String()}
	schemas := ast.Schemas{
		&ast.Schema{
			Package: "pkg",
			Objects: testutils.ObjectsMap(originalObject),
		},
	}
	originalBuilders := ast.Builders{
		{
			For:     originalObject,
			Package: "pkg",
			Name:    "Dashboard",
			Options: []ast.Option{
				{
					Name: "name",
					Args: []ast.Argument{argument},
					Assignments: []ast.Assignment{
						ast.ArgumentAssignment(ast.PathFromStructField(originalObject.Type.Struct.Fields[0]), argument),
					},
				},
			},
		},
	}

	rule := Duplicate(ByName("pkg", "Dashboard"), "NewDashboard", nil)
	updatedBuilders, err := rule(schemas, originalBuilders)
	req.NoError(err)

	req.Len(updatedBuilders, 2)
	req.Equal(originalBuilders[0], updatedBuilders[0])

	req.Equal("NewDashboard", updatedBuilders[1].Name)
	req.Equal([]string{"Duplicate[pkg.Dashboard]"}, updatedBuilders[1].VeneerTrail)
}

func TestInitialize(t *testing.T) {
	req := require.New(t)

	originalObject := ast.NewObject("pkg", "Dashboard", ast.NewStruct(
		ast.NewStructField("name", ast.String()),
	))
	argument := ast.Argument{Name: "name", Type: ast.String()}
	schemas := ast.Schemas{
		&ast.Schema{
			Package: "pkg",
			Objects: testutils.ObjectsMap(originalObject),
		},
	}
	originalBuilders := ast.Builders{
		{
			For:     originalObject,
			Package: "pkg",
			Name:    "Dashboard",
			Options: []ast.Option{
				{
					Name: "name",
					Args: []ast.Argument{argument},
					Assignments: []ast.Assignment{
						ast.ArgumentAssignment(ast.PathFromStructField(originalObject.Type.Struct.Fields[0]), argument),
					},
				},
			},
		},
	}

	rule := Initialize(
		ByName("pkg", "Dashboard"),
		[]Initialization{
			{PropertyPath: "name", Value: "great name, isn't it?"},
		},
	)
	updatedBuilders, err := rule(schemas, originalBuilders)
	req.NoError(err)

	expectedAssignments := []ast.Assignment{
		{
			Path:   ast.Path{{Identifier: "name", Type: ast.String()}},
			Value:  ast.AssignmentValue{Constant: "great name, isn't it?"},
			Method: ast.DirectAssignment,
		},
	}

	req.Len(updatedBuilders, 1)
	req.Equal(expectedAssignments, updatedBuilders[0].Constructor.Assignments)
}

func TestPromoteOptionsToConstructor(t *testing.T) {
	req := require.New(t)

	originalObject := ast.NewObject("pkg", "Dashboard", ast.NewStruct(
		ast.NewStructField("uid", ast.String()),
		ast.NewStructField("name", ast.String()),
	))
	argument := ast.Argument{Name: "name", Type: ast.String()}
	assignment := ast.ArgumentAssignment(ast.PathFromStructField(originalObject.Type.Struct.Fields[0]), argument)
	schemas := ast.Schemas{
		&ast.Schema{
			Package: "pkg",
			Objects: testutils.ObjectsMap(originalObject),
		},
	}
	originalBuilders := ast.Builders{
		{
			For:     originalObject,
			Package: "pkg",
			Name:    "Dashboard",
			Options: []ast.Option{
				{
					Name:        "name",
					Args:        []ast.Argument{argument},
					Assignments: []ast.Assignment{assignment},
				},
			},
		},
	}

	rule := PromoteOptionsToConstructor(
		ByName("pkg", "Dashboard"),
		[]string{"name"},
	)
	updatedBuilders, err := rule(schemas, originalBuilders)
	req.NoError(err)

	expectedArgs := []ast.Argument{argument}
	expectedAssignments := []ast.Assignment{assignment}

	req.Len(updatedBuilders, 1)
	req.Equal(expectedArgs, updatedBuilders[0].Constructor.Args)
	req.Equal(expectedAssignments, updatedBuilders[0].Constructor.Assignments)
}
