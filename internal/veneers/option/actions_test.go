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
		Comments: []string{"some comments"},
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
