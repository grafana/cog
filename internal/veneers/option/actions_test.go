package option

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestRenameAction(t *testing.T) {
	req := require.New(t)

	option := ast.Option{Name: "Name"}
	modifiedOpts := RenameAction("NewName")(ast.Builder{}, option)

	req.Len(modifiedOpts, 1)
	req.Equal("NewName", modifiedOpts[0].Name)
}

func TestOmitAction(t *testing.T) {
	req := require.New(t)

	option := ast.Option{Name: "Name"}
	modifiedOpts := OmitAction()(ast.Builder{}, option)

	req.Empty(modifiedOpts)
}

func TestPromoteToConstructorAction(t *testing.T) {
	req := require.New(t)

	option := ast.Option{Name: "Name", IsConstructorArg: false}
	modifiedOpts := PromoteToConstructorAction()(ast.Builder{}, option)

	req.Len(modifiedOpts, 1)
	req.True(modifiedOpts[0].IsConstructorArg)
}

func TestUnfoldBooleanAction(t *testing.T) {
	req := require.New(t)

	option := ast.Option{
		Args: []ast.Argument{
			{Name: "editable", Type: ast.Bool()},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "editable", Type: ast.Bool()},
			}, ast.Argument{Name: "editable", Type: ast.Bool()}),
		},
	}
	modifiedOpts := UnfoldBooleanAction(BooleanUnfold{
		OptionTrue:  "Editable",
		OptionFalse: "ReadOnly",
	})(ast.Builder{}, option)

	req.Len(modifiedOpts, 2)

	editableOpt := modifiedOpts[0]
	readonlyOpt := modifiedOpts[1]

	req.Equal(editableOpt.Name, "Editable")
	req.Len(editableOpt.Assignments, 1)
	req.Len(editableOpt.Args, 0)
	req.Equal(editableOpt.Assignments[0].Path.String(), "editable")
	req.Equal(editableOpt.Assignments[0].Value.Constant, true)

	req.Equal(readonlyOpt.Name, "ReadOnly")
	req.Len(readonlyOpt.Args, 0)
	req.Len(readonlyOpt.Assignments, 1)
	req.Equal(readonlyOpt.Assignments[0].Path.String(), "editable")
	req.Equal(readonlyOpt.Assignments[0].Value.Constant, false)
}

func TestDisjunctionAsOptionsAction_withDisjunction(t *testing.T) {
	req := require.New(t)

	disjunctionType := ast.NewDisjunction(ast.Types{
		ast.NewRef("dashboard", "Panel"),
		ast.NewRef("dashboard", "Row"),
	})

	option := ast.Option{
		Name: "Panel",
		Args: []ast.Argument{
			{Name: "panel", Type: disjunctionType},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "panel", Type: disjunctionType},
			}, ast.Argument{Name: "tags", Type: disjunctionType}),
		},
	}
	modifiedOpts := DisjunctionAsOptionsAction()(ast.Builder{}, option)

	req.Len(modifiedOpts, 2)

	req.Equal("panel", modifiedOpts[0].Name)
	req.Len(modifiedOpts[0].Args, 1)
	req.Equal(disjunctionType.Disjunction.Branches[0], modifiedOpts[0].Args[0].Type)
	req.Equal("panel", modifiedOpts[0].Args[0].Name)

	req.Equal("row", modifiedOpts[1].Name)
	req.Len(modifiedOpts[1].Args, 1)
	req.Equal(disjunctionType.Disjunction.Branches[1], modifiedOpts[1].Args[0].Type)
	req.Equal("row", modifiedOpts[1].Args[0].Name)
}

func TestDisjunctionAsOptionsAction_withDisjunctionStruct(t *testing.T) {
	req := require.New(t)

	panelType := ast.NewStruct()
	rowType := ast.NewStruct()
	panelOrRow := ast.NewStruct(
		ast.NewStructField("Panel", ast.NewRef("dashboard", "Panel")),
		ast.NewStructField("Row", ast.NewRef("dashboard", "Row")),
	)
	panelOrRow.Hints[ast.HintDiscriminatedDisjunctionOfRefs] = "not nil"
	ref := ast.NewRef("dashboard", "PanelOrRow")
	schema := &ast.Schema{
		Package: "dashboard",
		Objects: testutils.ObjectsMap(
			ast.NewObject("dashboard", "PanelOrRow", panelOrRow),
			ast.NewObject("dashboard", "Row", rowType),
			ast.NewObject("dashboard", "Panel", panelType),
		),
	}
	builder := ast.Builder{Schema: schema}

	option := ast.Option{
		Name: "Panel",
		Args: []ast.Argument{
			{Name: "panel", Type: ref},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "panel", Type: ref},
			}, ast.Argument{Name: "tags", Type: ref}),
		},
	}
	modifiedOpts := DisjunctionAsOptionsAction()(builder, option)

	req.Len(modifiedOpts, 2)

	req.Equal("Panel", modifiedOpts[0].Name)
	req.Len(modifiedOpts[0].Args, 1)
	req.Equal(ast.NewRef("dashboard", "Panel"), modifiedOpts[0].Args[0].Type)
	req.Equal("Panel", modifiedOpts[0].Args[0].Name)

	req.Equal("Row", modifiedOpts[1].Name)
	req.Len(modifiedOpts[1].Args, 1)
	req.Equal(ast.NewRef("dashboard", "Row"), modifiedOpts[1].Args[0].Type)
	req.Equal("Row", modifiedOpts[1].Args[0].Name)
}

func TestStructFieldsAsOptionsAction_withRefArg(t *testing.T) {
	req := require.New(t)

	timeType := ast.NewStruct(
		ast.NewStructField("from", ast.String()),
		ast.NewStructField("to", ast.String()),
		ast.NewStructField("auto", ast.Bool()),
	)
	ref := ast.NewRef("dashboard", "Time")
	schema := &ast.Schema{
		Package: "dashboard",
		Objects: testutils.ObjectsMap(
			ast.NewObject("dashboard", "Time", timeType),
		),
	}
	builder := ast.Builder{Schema: schema}

	option := ast.Option{
		Name: "Time",
		Args: []ast.Argument{
			{Name: "time", Type: ref},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "time", Type: ref},
			}, ast.Argument{Name: "editable", Type: ref}),
		},
	}
	modifiedOpts := StructFieldsAsOptionsAction("from", "to")(builder, option)

	req.Len(modifiedOpts, 2)

	req.Equal("from", modifiedOpts[0].Name)
	req.Len(modifiedOpts[0].Args, 1)
	req.Equal("from", modifiedOpts[0].Args[0].Name)
	req.Equal(ast.String(), modifiedOpts[0].Args[0].Type)
	req.Len(modifiedOpts[0].Assignments, 1)
	req.Equal("time.from", modifiedOpts[0].Assignments[0].Path.String())

	req.Equal("to", modifiedOpts[1].Name)
	req.Len(modifiedOpts[1].Args, 1)
	req.Equal("to", modifiedOpts[1].Args[0].Name)
	req.Equal(ast.String(), modifiedOpts[1].Args[0].Type)
	req.Len(modifiedOpts[1].Assignments, 1)
	req.Equal("time.to", modifiedOpts[1].Assignments[0].Path.String())
}

func TestArrayToAppendAction_withNoArgument(t *testing.T) {
	req := require.New(t)

	option := ast.Option{
		Assignments: []ast.Assignment{
			ast.ConstantAssignment(ast.Path{
				{Identifier: "editable", Type: ast.Bool()},
			}, true),
		},
	}
	modifiedOpts := ArrayToAppendAction()(ast.Builder{}, option)

	req.Equal([]ast.Option{option}, modifiedOpts)
}

func TestArrayToAppendAction_withNonArrayArgument(t *testing.T) {
	req := require.New(t)

	option := ast.Option{
		Args: []ast.Argument{
			{Name: "editable", Type: ast.Bool()},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "editable", Type: ast.Bool()},
			}, ast.Argument{Name: "editable", Type: ast.Bool()}),
		},
	}
	modifiedOpts := ArrayToAppendAction()(ast.Builder{}, option)

	req.Equal([]ast.Option{option}, modifiedOpts)
}

func TestArrayToAppendAction_withArrayArgument(t *testing.T) {
	req := require.New(t)

	// input
	option := ast.Option{
		Args: []ast.Argument{
			{Name: "tags", Type: ast.NewArray(ast.String())},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "tags", Type: ast.NewArray(ast.String())},
			}, ast.Argument{Name: "tags", Type: ast.NewArray(ast.String())}),
		},
	}

	// expected output
	expectedOption := ast.Option{
		Args: []ast.Argument{
			{Name: "tags", Type: ast.String()},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(
				ast.Path{
					{Identifier: "tags", Type: ast.NewArray(ast.String())},
				},
				ast.Argument{Name: "tags", Type: ast.String()},
				ast.Method(ast.AppendAssignment),
			),
		},
		VeneerTrail: []string{"ArrayToAppend"},
	}

	modifiedOpts := ArrayToAppendAction()(ast.Builder{}, option)

	req.Equal([]ast.Option{expectedOption}, modifiedOpts)
}

func TestStructFieldsAsArgumentsAction_withNoArgument(t *testing.T) {
	req := require.New(t)

	option := ast.Option{
		Assignments: []ast.Assignment{
			ast.ConstantAssignment(ast.Path{
				{Identifier: "editable", Type: ast.Bool()},
			}, true),
		},
	}
	modifiedOpts := StructFieldsAsArgumentsAction()(ast.Builder{}, option)

	req.Equal([]ast.Option{option}, modifiedOpts)
}

func TestStructFieldsAsArgumentsAction_withNonStructArgument(t *testing.T) {
	req := require.New(t)

	option := ast.Option{
		Args: []ast.Argument{
			{Name: "tags", Type: ast.NewArray(ast.String())},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "tags", Type: ast.NewArray(ast.String())},
			}, ast.Argument{Name: "tags", Type: ast.NewArray(ast.String())}),
		},
	}
	modifiedOpts := StructFieldsAsArgumentsAction()(ast.Builder{}, option)

	req.Equal([]ast.Option{option}, modifiedOpts)
}

func TestStructFieldsAsArgumentsAction_withStructArgument(t *testing.T) {
	req := require.New(t)

	structType := ast.NewStruct(
		ast.NewStructField("from", ast.String()),
		ast.NewStructField("to", ast.String()),
		ast.NewStructField("type", ast.String(ast.Value("time"))),
	)

	// input
	option := ast.Option{
		Args: []ast.Argument{
			{Name: "time", Type: structType},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "time", Type: structType},
			}, ast.Argument{Name: "time", Type: structType}),
		},
	}

	// expected
	expectedOption := ast.Option{
		Args: []ast.Argument{
			{Name: "from", Type: ast.String()},
			{Name: "to", Type: ast.String()},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "time", Type: structType},
				{Identifier: "from", Type: ast.String()},
			}, ast.Argument{Name: "from", Type: ast.String()}),
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "time", Type: structType},
				{Identifier: "to", Type: ast.String()},
			}, ast.Argument{Name: "to", Type: ast.String()}),
			ast.ConstantAssignment(ast.Path{
				{Identifier: "time", Type: structType},
				{Identifier: "type", Type: ast.String(ast.Value("time"))},
			}, "time"),
		},
		VeneerTrail: []string{"StructFieldsAsArguments"},
	}

	modifiedOpts := StructFieldsAsArgumentsAction()(ast.Builder{}, option)

	req.Equal([]ast.Option{expectedOption}, modifiedOpts)
}

func TestStructFieldsAsArgumentsAction_withArrayOfStructArgument(t *testing.T) {
	req := require.New(t)

	structType := ast.NewStruct(
		ast.NewStructField("from", ast.String()),
		ast.NewStructField("to", ast.String()),
		ast.NewStructField("type", ast.String(ast.Value("time"))),
	)

	// input
	option := ast.Option{
		Args: []ast.Argument{
			{Name: "time", Type: structType},
		},
		Assignments: []ast.Assignment{
			ast.ArgumentAssignment(ast.Path{
				{Identifier: "time", Type: ast.NewArray(structType)},
			}, ast.Argument{Name: "time", Type: structType}),
		},
	}

	// expected
	expectedOption := ast.Option{
		Args: []ast.Argument{
			{Name: "from", Type: ast.String()},
			{Name: "to", Type: ast.String()},
		},
		Assignments: []ast.Assignment{
			{
				Method: ast.AppendAssignment,
				Path:   ast.Path{{Identifier: "time", Type: ast.NewArray(structType)}},
				Value: ast.AssignmentValue{
					Envelope: &ast.AssignmentEnvelope{
						Type: structType,
						Values: []ast.EnvelopeFieldValue{
							{
								Path: ast.Path{{Identifier: "from", Type: ast.String()}},
								Value: ast.AssignmentValue{Argument: &ast.Argument{
									Name: "from",
									Type: ast.String(),
								}},
							},
							{
								Path: ast.Path{{Identifier: "to", Type: ast.String()}},
								Value: ast.AssignmentValue{Argument: &ast.Argument{
									Name: "to",
									Type: ast.String(),
								}},
							},
							{
								Path:  ast.Path{{Identifier: "type", Type: ast.String(ast.Value("time"))}},
								Value: ast.AssignmentValue{Constant: "time"},
							},
						},
					},
				},
			},
		},
		VeneerTrail: []string{"StructFieldsAsArguments"},
	}

	modifiedOpts := StructFieldsAsArgumentsAction()(ast.Builder{}, option)

	req.Equal([]ast.Option{expectedOption}, modifiedOpts)
}
