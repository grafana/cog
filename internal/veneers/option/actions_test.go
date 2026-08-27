package option

import (
	"testing"

	"github.com/grafana/cog/internal/logs"
	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
	"github.com/stretchr/testify/require"
)

func TestRenameAction(t *testing.T) {
	req := require.New(t)

	option := ir.Option{Name: "Name"}
	modifiedOpts, err := RenameAction("NewName")(RuleCtx{}, ir.Builder{}, option)
	req.NoError(err)

	req.Len(modifiedOpts, 1)
	req.Equal("NewName", modifiedOpts[0].Name)
}

func TestOmitAction(t *testing.T) {
	req := require.New(t)

	option := ir.Option{Name: "Name"}
	modifiedOpts, err := OmitAction()(RuleCtx{}, ir.Builder{}, option)
	req.NoError(err)

	req.Empty(modifiedOpts)
}

func TestUnfoldBooleanAction(t *testing.T) {
	req := require.New(t)

	option := ir.Option{
		Args: []ir.Argument{
			{Name: "editable", Type: ir.Bool()},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "editable", Type: ir.Bool()},
			}, ir.Argument{Name: "editable", Type: ir.Bool()}),
		},
	}
	modifiedOpts, err := UnfoldBooleanAction(BooleanUnfold{
		OptionTrue:  "Editable",
		OptionFalse: "ReadOnly",
	})(RuleCtx{}, ir.Builder{}, option)
	req.NoError(err)

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

func TestUnfoldBooleanAction_onNonBooleanDoesNothing(t *testing.T) {
	req := require.New(t)

	option := ir.Option{
		Args: []ir.Argument{
			{Name: "tags", Type: ir.NewArray(ir.String())},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "tags", Type: ir.NewArray(ir.String())},
			}, ir.Argument{Name: "tags", Type: ir.NewArray(ir.String())}),
		},
	}
	modifiedOpts, err := UnfoldBooleanAction(BooleanUnfold{
		OptionTrue:  "TrueOpt",
		OptionFalse: "FalseOpt",
	})(RuleCtx{Logger: logs.NoopLogger()}, ir.Builder{}, option)
	req.NoError(err)

	req.Len(modifiedOpts, 1)
	req.Equal(option, modifiedOpts[0])
}

func TestDisjunctionAsOptionsAction_withDisjunction(t *testing.T) {
	req := require.New(t)

	disjunctionType := ir.NewDisjunction(ir.Types{
		ir.NewRef("dashboard", "Panel"),
		ir.NewRef("dashboard", "Row"),
	})

	option := ir.Option{
		Name: "Panel",
		Args: []ir.Argument{
			{Name: "panel", Type: disjunctionType},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "panel", Type: disjunctionType},
			}, ir.Argument{Name: "tags", Type: disjunctionType}),
		},
	}
	modifiedOpts, err := DisjunctionAsOptionsAction(0)(RuleCtx{}, ir.Builder{}, option)
	req.NoError(err)

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

func TestDisjunctionAsOptionsAction_withDisjunctionAsSecondArg(t *testing.T) {
	req := require.New(t)

	disjunctionType := ir.NewDisjunction(ir.Types{
		ir.NewRef("dashboard", "Panel"),
		ir.NewRef("dashboard", "Row"),
	})

	option := ir.Option{
		Name: "Panel",
		Args: []ir.Argument{
			{Name: "key", Type: ir.String()},
			{Name: "panel", Type: disjunctionType},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{ // This assignment doesn't make sense, but for the purpose of this test it doesn't matter.
				{Identifier: "key", Type: ir.String()},
			}, ir.Argument{Name: "key", Type: ir.String()}),
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "panel", Type: disjunctionType},
			}, ir.Argument{Name: "tags", Type: disjunctionType}),
		},
	}
	modifiedOpts, err := DisjunctionAsOptionsAction(1)(RuleCtx{}, ir.Builder{}, option)
	req.NoError(err)

	req.Len(modifiedOpts, 2)

	req.Equal("panel", modifiedOpts[0].Name)
	req.Len(modifiedOpts[0].Args, 2)
	req.Len(modifiedOpts[0].Assignments, 2)
	req.Equal("key", modifiedOpts[0].Args[0].Name)
	req.Equal(ir.String(), modifiedOpts[0].Args[0].Type)
	req.Equal("panel", modifiedOpts[0].Args[1].Name)
	req.Equal(disjunctionType.Disjunction.Branches[0], modifiedOpts[0].Args[1].Type)

	req.Equal("row", modifiedOpts[1].Name)
	req.Len(modifiedOpts[1].Args, 2)
	req.Len(modifiedOpts[1].Assignments, 2)
	req.Equal("key", modifiedOpts[1].Args[0].Name)
	req.Equal(ir.String(), modifiedOpts[1].Args[0].Type)
	req.Equal("row", modifiedOpts[1].Args[1].Name)
	req.Equal(disjunctionType.Disjunction.Branches[1], modifiedOpts[1].Args[1].Type)
}

func TestDisjunctionAsOptionsAction_withDisjunctionStruct(t *testing.T) {
	req := require.New(t)

	panelType := ir.NewStruct()
	rowType := ir.NewStruct()
	panelOrRow := ir.NewStruct(
		ir.NewStructField("Panel", ir.NewRef("dashboard", "Panel")),
		ir.NewStructField("Row", ir.NewRef("dashboard", "Row")),
	)
	panelOrRow.Hints[ir.HintDiscriminatedDisjunctionOfRefs] = "not nil"
	ref := ir.NewRef("dashboard", "PanelOrRow")
	schema := &ir.Schema{
		Package: "dashboard",
		Objects: testutils.ObjectsMap(
			ir.NewObject("dashboard", "PanelOrRow", panelOrRow),
			ir.NewObject("dashboard", "Row", rowType),
			ir.NewObject("dashboard", "Panel", panelType),
		),
	}
	option := ir.Option{
		Name: "Panel",
		Args: []ir.Argument{
			{Name: "panel", Type: ref},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "panel", Type: ref},
			}, ir.Argument{Name: "tags", Type: ref}),
		},
	}
	ctx := RuleCtx{
		Schemas: ir.Schemas{schema},
	}
	modifiedOpts, err := DisjunctionAsOptionsAction(0)(ctx, ir.Builder{}, option)
	req.NoError(err)

	req.Len(modifiedOpts, 2)

	req.Equal("Panel", modifiedOpts[0].Name)
	req.Len(modifiedOpts[0].Args, 1)
	req.Equal(ir.NewRef("dashboard", "Panel"), modifiedOpts[0].Args[0].Type)
	req.Equal("Panel", modifiedOpts[0].Args[0].Name)

	req.Equal("Row", modifiedOpts[1].Name)
	req.Len(modifiedOpts[1].Args, 1)
	req.Equal(ir.NewRef("dashboard", "Row"), modifiedOpts[1].Args[0].Type)
	req.Equal("Row", modifiedOpts[1].Args[0].Name)
}

func TestDisjunctionAsOptionsAction_withDisjunctionStructAsSecondArg(t *testing.T) {
	req := require.New(t)

	panelType := ir.NewStruct()
	rowType := ir.NewStruct()
	panelOrRow := ir.NewStruct(
		ir.NewStructField("Panel", ir.NewRef("dashboard", "Panel")),
		ir.NewStructField("Row", ir.NewRef("dashboard", "Row")),
	)
	panelOrRow.Hints[ir.HintDiscriminatedDisjunctionOfRefs] = "not nil"
	ref := ir.NewRef("dashboard", "PanelOrRow")
	schema := &ir.Schema{
		Package: "dashboard",
		Objects: testutils.ObjectsMap(
			ir.NewObject("dashboard", "PanelOrRow", panelOrRow),
			ir.NewObject("dashboard", "Row", rowType),
			ir.NewObject("dashboard", "Panel", panelType),
		),
	}
	option := ir.Option{
		Name: "Panel",
		Args: []ir.Argument{
			{Name: "key", Type: ir.String()},
			{Name: "panel", Type: ref},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{ // This assignment doesn't make sense, but for the purpose of this test it doesn't matter.
				{Identifier: "key", Type: ir.String()},
			}, ir.Argument{Name: "key", Type: ir.String()}),
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "panel", Type: ref},
			}, ir.Argument{Name: "tags", Type: ref}),
		},
	}
	ctx := RuleCtx{
		Schemas: ir.Schemas{schema},
	}
	modifiedOpts, err := DisjunctionAsOptionsAction(1)(ctx, ir.Builder{}, option)
	req.NoError(err)

	req.Len(modifiedOpts, 2)

	req.Equal("Panel", modifiedOpts[0].Name)
	req.Len(modifiedOpts[0].Args, 2)
	req.Len(modifiedOpts[0].Assignments, 2)
	req.Equal("key", modifiedOpts[0].Args[0].Name)
	req.Equal(ir.String(), modifiedOpts[0].Args[0].Type)
	req.Equal(ir.NewRef("dashboard", "Panel"), modifiedOpts[0].Args[1].Type)
	req.Equal("Panel", modifiedOpts[0].Args[1].Name)

	req.Equal("Row", modifiedOpts[1].Name)
	req.Len(modifiedOpts[1].Args, 2)
	req.Len(modifiedOpts[1].Assignments, 2)
	req.Equal("key", modifiedOpts[0].Args[0].Name)
	req.Equal(ir.String(), modifiedOpts[0].Args[0].Type)
	req.Equal(ir.NewRef("dashboard", "Row"), modifiedOpts[1].Args[1].Type)
	req.Equal("Row", modifiedOpts[1].Args[1].Name)
}

func TestStructFieldsAsOptionsAction_withRefArg(t *testing.T) {
	req := require.New(t)

	timeType := ir.NewStruct(
		ir.NewStructField("from", ir.String()),
		ir.NewStructField("to", ir.String()),
		ir.NewStructField("auto", ir.Bool()),
	)
	ref := ir.NewRef("dashboard", "Time")
	schema := &ir.Schema{
		Package: "dashboard",
		Objects: testutils.ObjectsMap(
			ir.NewObject("dashboard", "Time", timeType),
		),
	}
	option := ir.Option{
		Name: "Time",
		Args: []ir.Argument{
			{Name: "time", Type: ref},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "time", Type: ref},
			}, ir.Argument{Name: "editable", Type: ref}),
		},
	}
	ctx := RuleCtx{
		Schemas: ir.Schemas{schema},
	}
	modifiedOpts, err := StructFieldsAsOptionsAction("from", "to")(ctx, ir.Builder{}, option)
	req.NoError(err)

	req.Len(modifiedOpts, 2)

	req.Equal("from", modifiedOpts[0].Name)
	req.Len(modifiedOpts[0].Args, 1)
	req.Equal("from", modifiedOpts[0].Args[0].Name)
	req.Equal(ir.String(), modifiedOpts[0].Args[0].Type)
	req.Len(modifiedOpts[0].Assignments, 1)
	req.Equal("time.from", modifiedOpts[0].Assignments[0].Path.String())

	req.Equal("to", modifiedOpts[1].Name)
	req.Len(modifiedOpts[1].Args, 1)
	req.Equal("to", modifiedOpts[1].Args[0].Name)
	req.Equal(ir.String(), modifiedOpts[1].Args[0].Type)
	req.Len(modifiedOpts[1].Assignments, 1)
	req.Equal("time.to", modifiedOpts[1].Assignments[0].Path.String())
}

func TestArrayToAppendAction_withNoArgument(t *testing.T) {
	req := require.New(t)

	option := ir.Option{
		Assignments: []ir.Assignment{
			ir.ConstantAssignment(ir.Path{
				{Identifier: "editable", Type: ir.Bool()},
			}, true),
		},
	}
	modifiedOpts, err := ArrayToAppendAction()(RuleCtx{Logger: logs.NoopLogger()}, ir.Builder{}, option)
	req.NoError(err)

	req.Equal([]ir.Option{option}, modifiedOpts)
}

func TestArrayToAppendAction_withNonArrayArgument(t *testing.T) {
	req := require.New(t)

	option := ir.Option{
		Args: []ir.Argument{
			{Name: "editable", Type: ir.Bool()},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "editable", Type: ir.Bool()},
			}, ir.Argument{Name: "editable", Type: ir.Bool()}),
		},
	}
	modifiedOpts, err := ArrayToAppendAction()(RuleCtx{Logger: logs.NoopLogger()}, ir.Builder{}, option)
	req.NoError(err)

	req.Equal([]ir.Option{option}, modifiedOpts)
}

func TestArrayToAppendAction_withArrayArgument(t *testing.T) {
	req := require.New(t)

	// input
	option := ir.Option{
		Args: []ir.Argument{
			{Name: "tags", Type: ir.NewArray(ir.String())},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "tags", Type: ir.NewArray(ir.String())},
			}, ir.Argument{Name: "tags", Type: ir.NewArray(ir.String())}),
		},
	}

	// expected output
	expectedOption := ir.Option{
		Args: []ir.Argument{
			{Name: "tag", Type: ir.String()},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(
				ir.Path{
					{Identifier: "tags", Type: ir.NewArray(ir.String())},
				},
				ir.Argument{Name: "tag", Type: ir.String()},
				ir.Method(ir.AppendAssignment),
			),
		},
		VeneerTrail: []string{"ArrayToAppend"},
	}

	modifiedOpts, err := ArrayToAppendAction()(RuleCtx{}, ir.Builder{}, option)
	req.NoError(err)

	req.Equal([]ir.Option{expectedOption}, modifiedOpts)
}

func TestStructFieldsAsArgumentsAction_withNoArgument(t *testing.T) {
	req := require.New(t)

	option := ir.Option{
		Assignments: []ir.Assignment{
			ir.ConstantAssignment(ir.Path{
				{Identifier: "editable", Type: ir.Bool()},
			}, true),
		},
	}
	modifiedOpts, err := StructFieldsAsArgumentsAction()(RuleCtx{Logger: logs.NoopLogger()}, ir.Builder{}, option)
	req.NoError(err)

	req.Equal([]ir.Option{option}, modifiedOpts)
}

func TestStructFieldsAsArgumentsAction_withNonStructArgument(t *testing.T) {
	req := require.New(t)

	option := ir.Option{
		Args: []ir.Argument{
			{Name: "tags", Type: ir.NewArray(ir.String())},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "tags", Type: ir.NewArray(ir.String())},
			}, ir.Argument{Name: "tags", Type: ir.NewArray(ir.String())}),
		},
	}
	modifiedOpts, err := StructFieldsAsArgumentsAction()(RuleCtx{Logger: logs.NoopLogger()}, ir.Builder{}, option)
	req.NoError(err)

	req.Equal([]ir.Option{option}, modifiedOpts)
}

func TestStructFieldsAsArgumentsAction_withStructArgument(t *testing.T) {
	req := require.New(t)

	structType := ir.NewStruct(
		ir.NewStructField("from", ir.String()),
		ir.NewStructField("to", ir.String()),
		ir.NewStructField("type", ir.String(ir.Value("time"))),
	)

	// input
	option := ir.Option{
		Args: []ir.Argument{
			{Name: "time", Type: structType},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "time", Type: structType},
			}, ir.Argument{Name: "time", Type: structType}),
		},
	}

	// expected
	expectedOption := ir.Option{
		Args: []ir.Argument{
			{Name: "from", Type: ir.String()},
			{Name: "to", Type: ir.String()},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "time", Type: structType},
				{Identifier: "from", Type: ir.String()},
			}, ir.Argument{Name: "from", Type: ir.String()}),
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "time", Type: structType},
				{Identifier: "to", Type: ir.String()},
			}, ir.Argument{Name: "to", Type: ir.String()}),
			ir.ConstantAssignment(ir.Path{
				{Identifier: "time", Type: structType},
				{Identifier: "type", Type: ir.String(ir.Value("time"))},
			}, "time"),
		},
		VeneerTrail: []string{"StructFieldsAsArguments"},
	}

	modifiedOpts, err := StructFieldsAsArgumentsAction()(RuleCtx{}, ir.Builder{}, option)
	req.NoError(err)

	req.Equal([]ir.Option{expectedOption}, modifiedOpts)
}

func TestStructFieldsAsArgumentsAction_withArrayOfStructArgument(t *testing.T) {
	req := require.New(t)

	structType := ir.NewStruct(
		ir.NewStructField("from", ir.String()),
		ir.NewStructField("to", ir.String()),
		ir.NewStructField("type", ir.String(ir.Value("time"))),
	)

	// input
	option := ir.Option{
		Args: []ir.Argument{
			{Name: "time", Type: structType},
		},
		Assignments: []ir.Assignment{
			ir.ArgumentAssignment(ir.Path{
				{Identifier: "time", Type: ir.NewArray(structType)},
			}, ir.Argument{Name: "time", Type: structType}),
		},
	}

	// expected
	expectedOption := ir.Option{
		Args: []ir.Argument{
			{Name: "from", Type: ir.String()},
			{Name: "to", Type: ir.String()},
		},
		Assignments: []ir.Assignment{
			{
				Method: ir.AppendAssignment,
				Path:   ir.Path{{Identifier: "time", Type: ir.NewArray(structType)}},
				Value: ir.AssignmentValue{
					Envelope: &ir.AssignmentEnvelope{
						Type: structType,
						Values: []ir.EnvelopeFieldValue{
							{
								Path: ir.Path{{Identifier: "from", Type: ir.String()}},
								Value: ir.AssignmentValue{Argument: &ir.Argument{
									Name: "from",
									Type: ir.String(),
								}},
							},
							{
								Path: ir.Path{{Identifier: "to", Type: ir.String()}},
								Value: ir.AssignmentValue{Argument: &ir.Argument{
									Name: "to",
									Type: ir.String(),
								}},
							},
							{
								Path:  ir.Path{{Identifier: "type", Type: ir.String(ir.Value("time"))}},
								Value: ir.AssignmentValue{Constant: "time"},
							},
						},
					},
				},
			},
		},
		VeneerTrail: []string{"StructFieldsAsArguments"},
	}

	modifiedOpts, err := StructFieldsAsArgumentsAction()(RuleCtx{}, ir.Builder{}, option)
	req.NoError(err)

	req.Equal([]ir.Option{expectedOption}, modifiedOpts)
}

func TestDuplicateAction(t *testing.T) {
	req := require.New(t)

	option := ir.Option{Name: "Name"}
	modifiedOpts, err := DuplicateAction("Duplicated")(RuleCtx{}, ir.Builder{}, option)
	req.NoError(err)

	req.Len(modifiedOpts, 2)
	req.Equal("Name", modifiedOpts[0].Name)
	req.Equal("Duplicated", modifiedOpts[1].Name)
}
