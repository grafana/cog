package option

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
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
		Assignments: []ast.Assignment{
			{
				Path: ast.Path{
					{Identifier: "editable", Type: ast.Bool()},
				},
			},
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
	req.Equal(editableOpt.Assignments[0].Path.String(), "editable")
	req.Equal(editableOpt.Assignments[0].Value.Constant, true)

	req.Equal(readonlyOpt.Name, "ReadOnly")
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
			{
				Path: ast.Path{
					{Identifier: "panel", Type: disjunctionType},
				},
			},
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
		Objects: []ast.Object{
			ast.NewObject("dashboard", "PanelOrRow", panelOrRow),
			ast.NewObject("dashboard", "Row", rowType),
			ast.NewObject("dashboard", "Panel", panelType),
		},
	}
	builder := ast.Builder{Schema: schema}

	option := ast.Option{
		Name: "Panel",
		Args: []ast.Argument{
			{Name: "panel", Type: ref},
		},
		Assignments: []ast.Assignment{
			{
				Path: ast.Path{
					{Identifier: "panel", Type: ref},
				},
			},
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