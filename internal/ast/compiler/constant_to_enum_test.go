package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestConstantToEnum(t *testing.T) {
	// Prepare test input
	strAlias := ast.NewObject("sandbox", "String", ast.String())
	strConstant := ast.NewObject("sandbox", "Mode", ast.String(ast.Value("auto")))
	notTargetedStrConstant := ast.NewObject("sandbox", "Alignment", ast.String(ast.Value("center")))
	intConstant := ast.NewObject("sandbox", "DefaultSize", ast.NewScalar(ast.KindInt32, ast.Value(42)))
	obj := ast.NewObject("sandbox", "Obj", ast.NewStruct(ast.NewStructField("foo", ast.String())))
	schema := &ast.Schema{
		Package: "sandbox",
		Objects: testutils.ObjectsMap(
			strAlias,
			strConstant,
			notTargetedStrConstant,
			intConstant,
			obj,
		),
	}

	newEnum := ast.NewObject("sandbox", "Mode", ast.NewEnum([]ast.EnumValue{
		{
			Type:  ast.String(),
			Name:  "auto",
			Value: "auto",
		},
	}))
	newEnum.AddToPassesTrail("ConstantToEnum")
	expected := &ast.Schema{
		Package: "sandbox",
		Objects: testutils.ObjectsMap(
			strAlias,
			newEnum,
			notTargetedStrConstant,
			intConstant,
			obj,
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &ConstantToEnum{
		Objects: []ObjectReference{
			{Package: "sandbox", Object: "String"},
			{Package: "sandbox", Object: "Mode"},
			{Package: "sandbox", Object: "DefaultSize"},
			{Package: "sandbox", Object: "Obj"},
		},
	}, schema, expected)
}
