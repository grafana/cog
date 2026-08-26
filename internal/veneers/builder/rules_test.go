package builder

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestDuplicate(t *testing.T) {
	req := require.New(t)

	originalObject := ir.NewObject("pkg", "Dashboard", ir.NewStruct(
		ir.NewStructField("name", ir.String()),
	))
	argument := ir.Argument{Name: "title", Type: ir.String()}
	schemas := ir.Schemas{
		&ir.Schema{
			Package: "pkg",
			Objects: testutils.ObjectsMap(originalObject),
		},
	}
	originalBuilders := ir.Builders{
		{
			For:     originalObject,
			Package: "pkg",
			Name:    "Dashboard",
			Options: []ir.Option{
				{
					Name: "name",
					Args: []ir.Argument{argument},
					Assignments: []ir.Assignment{
						ir.ArgumentAssignment(ir.PathFromStructField(originalObject.Type.Struct.Fields[0]), argument),
					},
				},
			},
		},
	}

	action := DuplicateAction("NewDashboard", nil)
	updatedBuilders, err := action(RuleCtx{Schemas: schemas, Builders: originalBuilders}, originalBuilders)
	req.NoError(err)

	req.Len(updatedBuilders, 2)
	req.Equal(originalBuilders[0], updatedBuilders[0])

	req.Equal("NewDashboard", updatedBuilders[1].Name)
	req.Equal([]string{"Duplicate[pkg.Dashboard]"}, updatedBuilders[1].VeneerTrail)
}

func TestInitialize(t *testing.T) {
	req := require.New(t)

	originalObject := ir.NewObject("pkg", "Dashboard", ir.NewStruct(
		ir.NewStructField("name", ir.String()),
	))
	argument := ir.Argument{Name: "name", Type: ir.String()}
	schemas := ir.Schemas{
		&ir.Schema{
			Package: "pkg",
			Objects: testutils.ObjectsMap(originalObject),
		},
	}
	originalBuilders := ir.Builders{
		{
			For:     originalObject,
			Package: "pkg",
			Name:    "Dashboard",
			Options: []ir.Option{
				{
					Name: "name",
					Args: []ir.Argument{argument},
					Assignments: []ir.Assignment{
						ir.ArgumentAssignment(ir.PathFromStructField(originalObject.Type.Struct.Fields[0]), argument),
					},
				},
			},
		},
	}

	action := InitializeAction([]Initialization{
		{PropertyPath: "name", Value: "great name, isn't it?"},
	})
	updatedBuilders, err := action(RuleCtx{Schemas: schemas, Builders: originalBuilders}, originalBuilders)
	req.NoError(err)

	expectedAssignments := []ir.Assignment{
		{
			Path:   ir.Path{{Identifier: "name", Type: ir.String()}},
			Value:  ir.AssignmentValue{Constant: "great name, isn't it?"},
			Method: ir.DirectAssignment,
		},
	}

	req.Len(updatedBuilders, 1)
	req.Equal(expectedAssignments, updatedBuilders[0].Constructor.Assignments)
}

func TestPromoteOptionsToConstructor(t *testing.T) {
	req := require.New(t)

	originalObject := ir.NewObject("pkg", "Dashboard", ir.NewStruct(
		ir.NewStructField("uid", ir.String()),
		ir.NewStructField("name", ir.String()),
	))
	argument := ir.Argument{Name: "name", Type: ir.String()}
	assignment := ir.ArgumentAssignment(ir.PathFromStructField(originalObject.Type.Struct.Fields[0]), argument)
	schemas := ir.Schemas{
		&ir.Schema{
			Package: "pkg",
			Objects: testutils.ObjectsMap(originalObject),
		},
	}
	originalBuilders := ir.Builders{
		{
			For:     originalObject,
			Package: "pkg",
			Name:    "Dashboard",
			Options: []ir.Option{
				{
					Name:        "name",
					Args:        []ir.Argument{argument},
					Assignments: []ir.Assignment{assignment},
				},
			},
		},
	}

	action := PromoteOptionsToConstructorAction([]string{"name"})
	updatedBuilders, err := action(RuleCtx{Schemas: schemas, Builders: originalBuilders}, originalBuilders)
	req.NoError(err)

	expectedArgs := []ir.Argument{argument}
	expectedAssignments := []ir.Assignment{assignment}

	req.Len(updatedBuilders, 1)
	req.Equal(expectedArgs, updatedBuilders[0].Constructor.Args)
	req.Equal(expectedAssignments, updatedBuilders[0].Constructor.Assignments)
}

func TestGenericAction(t *testing.T) {
	req := require.New(t)

	originalObject := ir.NewObject("generic", "Panel", ir.NewStruct(
		ir.NewStructField("uid", ir.String()),
		ir.NewStructField("name", ir.String()),
	))

	schemas := ir.Schemas{
		&ir.Schema{
			Package: "generic",
			Objects: testutils.ObjectsMap(originalObject),
		},
	}
	originalBuilders := ir.Builders{
		{
			For:     originalObject,
			Package: "generic",
		},
	}

	expectedBuilders := ir.Builders{
		{
			For:         originalObject,
			Package:     "generic",
			IsGeneric:   true,
			VeneerTrail: []string{"Generic[selector=generic.Panel]"},
		},
	}

	action := GenericAction(ByObjectName("generic", "Panel"))
	updatedBuilders, err := action(RuleCtx{Schemas: schemas, Builders: originalBuilders}, originalBuilders)
	req.NoError(err)

	req.Equal(expectedBuilders, updatedBuilders)
}
